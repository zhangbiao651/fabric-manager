package sdkInit

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/zhangbiao651/fabric-manager/sdkInit"
)

const (
	orgID     = "Org1MSP"
	userID    = "Admin"
	channelID = "mychannel"
)

func TestSetupSdk(t *testing.T) {
	path := "/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
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
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(peers)
}
