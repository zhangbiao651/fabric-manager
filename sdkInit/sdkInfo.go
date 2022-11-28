package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// 组织信息
type OrgInfo struct {
	OrgAdminUser          string                     // 组织的管理员用户
	OrgName               string                     // 组织名称
	OrgMspId              string                     // 组织的 MSPID
	OrgUser               string                     // 组织的普通用户
	OrgMspClient          *mspclient.Client          //组织的 MSP 客户端
	OrgAdminClientContext *contextAPI.ClientProvider // 组织管理员客户端的上下文
	OrgResMgmt            *resmgmt.Client            // 组织的更新客户端
	OrgPeerNum            int                        // 组织的节点个数
	OrgAnchorFile         string                     // 组织锚节点文件位置
	//	OrgConfigPath         string

	Peers map[string]*fabAPI.Peer // 组织下的节点
}

type SdkEnvInfo struct {
	// 通道信息
	ChannelID     string // like "simplecc"
	ChannelConfig string // like os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric-samples/test-network/channel-artifacts/testchannel.tx"

	// 组织信息
	OrgNames []string
	Orgs     map[string]*OrgInfo
	// 排序服务节点信息
	OrdererAdminUser     string // like "Admin"
	OrdererOrgName       string // like "OrdererOrg"
	OrdererEndpoint      string
	OrdererClientContext *contextAPI.ClientProvider
	// 链码信息
	ChaincodeID      string
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeVersion string
	ChClient         *channel.Client
	EvClient         *event.Client
}

type Application struct {
	SdkEnvInfo *SdkEnvInfo
}
