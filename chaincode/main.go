package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// ChainCode 链码结构体
type ChainCode struct {
}

// Init 链码初始化
func (c *ChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

//Invoke 链码调用
func (c *ChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "addRecord":
		return c.addRecord(stub, args)
	case "updateRecord":
		return c.updateRecord(stub, args)
	case "queryRecordByKey":
		return c.queryRecordByKey(stub, args)
	case "queryRecordByPatientID":
		return c.queryRecordByPatientID(stub, args)
	case "queryRecordByDoctorID":
		return c.queryRecordByDoctorID(stub, args)
	default:
		return shim.Error("调用方法不存在")
	}
}

func main() {
	err := shim.Start(new(ChainCode))
	if err != nil {
		fmt.Printf("启动链码时发生错误：%s", err)
	}
}
