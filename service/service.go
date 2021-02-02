package service

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

// 测试是否实现接口
var _ Service = &internal{}

// 实现服务的内部类
type internal struct {
	chaincodeID string
	client      *channel.Client
}

// Service 用于与底层链码交互，执行功能的接口，由内部internal类实现该接口
type Service interface {
	// TODO: 需实现的方法列表
}

// NewService 返回Service方便调用
func NewService(chaincodeID string, client *channel.Client) Service {
	return &internal{
		chaincodeID: chaincodeID,
		client:      client,
	}
}
