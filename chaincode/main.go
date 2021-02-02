package main

import (
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
	_, _ = stub.GetFunctionAndParameters()
	// TODO: 链码业务逻辑
	return shim.Success(nil)
}
