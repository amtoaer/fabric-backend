package sdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

//InitInfo 用于初始化SDK的必要信息
type InitInfo struct {
	ChannelID      string
	ChannelConfig  string
	OrgName        string
	OrgAdmin       string
	OrdererOrgName string
	OrgResMgmt     *resmgmt.Client

	ChaincodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
	UserName        string
}
