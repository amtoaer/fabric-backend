package main

import (
	"github.com/amtoaer/fabric-backend/sdk"
)

const (
	configFile  = "config.yaml"
	initialized = false
	// ChainCode 链码名
	ChainCode = "simplecc"
)

func main() {
	initInfo := &sdk.InitInfo{
		ChannelID:      "kevinkongyixueyuan",
		ChannelConfig:  "./fixtures/artifacts/channel.tx",
		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",
	}
	mySDK, err := sdk.SetupSDK(configFile, initialized)
	if err != nil {
		panic(err)
	}
	defer mySDK.Close()

	err = sdk.CreateChannel(mySDK, initInfo)

	if err != nil {
		panic(err)
	}

}
