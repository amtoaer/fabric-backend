package service

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

// Record 电子病历结构体（demo）
type Record struct {
	ObjectType  string
	PatientName string
	PatientID   string
	DoctorName  string
	DoctorID    string
	Content     string
	Historys    []HistoryItem
}

// HistoryItem 电子病历历史结构体
type HistoryItem struct {
	TxID    string
	History Record
}

// 实现服务的内部类
type internal struct {
	chaincodeID string
	client      *channel.Client
}

// Service 用于与底层链码交互，执行功能的接口，由内部internal类实现该接口
type Service interface {
	// 添加病历
	AddRecord(Record) (string, error)
	// 更新病历
	UpdateRecord(Record) (string, error)
	// 通过病人ID查找病历
	QueryRecordByPatientID(string) (string, error)
	// 通过医生ID查找病历
	QueryRecordByDoctorID(string) (string, error)
	// 通过医生和病人的ID（key）查找病历（带有修改历史）
	QueryRecordByKey(string, string) (string, error)
}

// NewService 返回Service方便调用
func NewService(chaincodeID string, client *channel.Client) Service {
	return &internal{
		chaincodeID: chaincodeID,
		client:      client,
	}
}
