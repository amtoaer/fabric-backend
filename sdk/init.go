package sdk

import (
	"fmt"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// SetupSDK 通过配置文件初始化SDK
func SetupSDK(configFile string, initialized bool) (*fabsdk.FabricSDK, error) {
	if initialized {
		return nil, fmt.Errorf("Fabric SDK已经实例化")
	}
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败：%v", err)
	}
	fmt.Println("Fabric SDK初始化成功")
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
	fmt.Println("通道已成功创建")
	info.OrgResMgmt = resMgmtClient
	err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("Peers加入通道失败： %v", err)
	}
	fmt.Println("peers 已成功加入通道.")
	return nil
}
