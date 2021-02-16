package main

import (
	"fmt"
	"os"
	"path"

	"github.com/amtoaer/fabric-backend/sdk"
	"github.com/amtoaer/fabric-backend/service"
	"github.com/amtoaer/fabric-backend/web"
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
	// 创建服务
	s := service.NewService(chainCode, channelClient)

	// 对服务进行部分测试
	s.AddRecord(service.Record{
		PatientName: "病人1",
		PatientID:   "1",
		DoctorName:  "医生1",
		DoctorID:    "1",
		Content:     "测试内容1",
	})
	s.AddRecord(service.Record{
		PatientName: "病人1",
		PatientID:   "1",
		DoctorName:  "医生2",
		DoctorID:    "2",
		Content:     "测试内容2",
	})
	s.AddRecord(service.Record{
		PatientName: "病人2",
		PatientID:   "2",
		DoctorName:  "医生1",
		DoctorID:    "1",
		Content:     "测试内容3",
	})
	result, err := s.QueryRecordByKey("1", "2")
	if err != nil {
		panic("第一次查询失败")
	}
	fmt.Printf("第一次查询结果：%v\n", result)
	result, err = s.QueryRecordByDoctorID("1")
	if err != nil {
		panic("第二次查询失败")
	}
	fmt.Printf("第二次查询结果：%v\n", result)
	result, err = s.QueryRecordByPatientID("1")
	if err != nil {
		panic("第三次查询失败")
	}
	fmt.Printf("第三次查询结果：%v\n", result)

	// 将服务注入到web包中
	web.SetService(s)
	// 拿到组装好的gin路由
	router := web.NewRouter()
	// 启用web服务并阻塞
	router.Run(":8000")
}
