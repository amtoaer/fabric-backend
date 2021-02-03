package main

import (
	"os"
	"path"

	"github.com/amtoaer/fabric-backend/sdk"
	"github.com/amtoaer/fabric-backend/service"
)

const (
	configFile = "config.yaml"
	// ChainCode 链码名
	chainCode = "simplecc"
)

func main() {
	initInfo := &sdk.InitInfo{
		ChannelID:      "kevinkongyixueyuan",
		ChannelConfig:  path.Join(os.Getenv("GOPATH"), "src/allwens.work/fabric-backend/fixtures/artifacts/channel.tx"),
		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID:     chainCode,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "allwens.work/fabric-backend/chaincode",
		UserName:        "User1",
	}
	// 初始化SDK
	mySDK, err := sdk.SetupSDK(configFile)
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
