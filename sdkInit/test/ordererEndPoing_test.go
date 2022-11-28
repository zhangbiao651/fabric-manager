package sdkInit

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

var (
	ordererEndpoint = "orderer.zhangbiao651.top"
)

func TestWithOrdererEndPoint(t *testing.T) {
	orderer := resmgmt.WithOrdererEndpoint(ordererEndpoint)
	fmt.Println(orderer)
}
