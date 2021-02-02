package main

import (
	"os"

	"github.com/amtoaer/fabric-backend/sdk"
	"github.com/amtoaer/fabric-backend/service"
)

const (
	configFile  = "config.yaml"
	initialized = false
	// ChainCode 链码名
	chainCode = "simplecc"
)

func main() {
	initInfo := &sdk.InitInfo{
		ChannelID:      "kevinkongyixueyuan",
		ChannelConfig:  "./fixtures/artifacts/channel.tx",
		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID: chainCode,
		// TODO: 因历史原因，fabric sdk创建链码包强制依赖GOPATH寻址，当前gomod模式需考虑迁移目录
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "allwens.work/fabric-backend/chaincode",
		UserName:        "User1",
	}
	// 初始化SDK
	mySDK, err := sdk.SetupSDK(configFile, initialized)
	if err != nil {
		panic(err)
	}
	defer mySDK.Close()
	// 创建通道并将peers加入通道
	err = sdk.CreateChannel(mySDK, initInfo)
	if err != nil {
		panic(err)
	}
	// 安装并实例化链码，拿到客户端
	channelClient, err := sdk.InstallAndInstantiateCC(mySDK, initInfo)
	if err != nil {
		panic(err)
	}
	// TODO: 该处创建的服务应该通过web中间层调用，待实现
	// 创建服务
	service.NewService(chainCode, channelClient)
}
