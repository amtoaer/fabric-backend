package sdk

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

const chaincodeVersion = "1.0.0"

// SetupSDK 通过配置文件初始化SDK
func SetupSDK(configFile string, initialized bool) (*fabsdk.FabricSDK, error) {
	if initialized {
		return nil, fmt.Errorf("Fabric SDK已经实例化")
	}
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败：%v", err)
	}
	fmt.Println("Fabric SDK初始化成功。")
	return sdk, nil
}

// CreateChannel 使用SDK创建通道
func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("创建资源管理客户端失败")
	}
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("创建通道管理客户端失败：%v", err)
	}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("创建 Org MSP 客户端失败：%v", err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名标识失败：%v", err)
	}
	req := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	_, err = resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败：%v", err)
	}
	fmt.Println("通道已成功创建。")
	info.OrgResMgmt = resMgmtClient
	err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("Peers加入通道失败： %v", err)
	}
	fmt.Println("peers 已成功加入通道。")
	return nil
}

// InstallAndInstantiateCC 安装并实例化链码
func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	fmt.Println("开始安装链码...")
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return nil, fmt.Errorf("创建链码包失败：%v", err)
	}
	req := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: chaincodeVersion, Package: ccPkg}
	_, err = info.OrgResMgmt.InstallCC(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("安装链码失败: %v", err)
	}
	fmt.Println("链码安装成功，开始实例化链码...")
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.kevin.kongyixueyuan.com"})
	instantiateCCReq := resmgmt.InstantiateCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: chaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy}
	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, fmt.Errorf("实例化链码失败: %v", err)
	}
	fmt.Println("链码实例化成功。")
	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// 返回客户端实例，用于查询链码、执行链码、注册/反注册特定通道的链码事件
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}
	fmt.Println("通道客户端创建成功。")
	return channelClient, nil
}
