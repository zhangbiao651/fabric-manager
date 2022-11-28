package sdkInit

import (
	"fmt"
	"testing"

	mb "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/policydsl"
	"github.com/zhangbiao651/fabric-manager/sdkInit"
)

var (
	packageID = "simplecc_1.0.0:eceb46fbfcc539c93faa3eed13820ceb67b2c401d48ae4ac8f36f933c7db862d"
	mspIDs    = []string{"avvv"}

	chaincodeID      = "mychannel"
	chaincodeVersion = "1.0.0"
	orgname          = "Org1MSP"
	path             = "/home/zb/fabric/project/v0.0.1/fabric-manager/config_test.yaml"
	orguser          = "Admin"
)

func TestLifecycleApprove(t *testing.T) {
	userID := "Admin"
	orgID := "Org1MSP"
	channelID := "mychannel"

	sdk, err := fabsdk.New(config.FromFile(path))
	if err != nil {
		fmt.Println(err)
	}
	defer sdk.Close()

	var options_user fabsdk.ContextOption
	var options_org fabsdk.ContextOption

	options_user = fabsdk.WithUser(userID)
	options_org = fabsdk.WithOrg(orgID)
	oac := sdk.Context(options_user, options_org)
	peers, err := sdkInit.DiscoverLocalPeers(oac, 1)

	ccPolicy := policydsl.SignedByNOutOfGivenRole(int32(len(mspIDs)), mb.MSPRole_MEMBER, mspIDs)
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              chaincodeID,
		Version:           chaincodeVersion,
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	cli, err := resmgmt.New(oac)
	c := *cli
	if _, err := c.LifecycleApproveCC(channelID, approveCCReq, resmgmt.WithTargets(peers...), resmgmt.WithOrdererEndpoint("orderer.zhangbiao651.top"), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		fmt.Println(err)
	}
}
