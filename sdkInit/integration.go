package sdkInit

import (
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/configtx"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// @title DiscoverLocaPeers
// @description  获取组织中的节点
// @author	zhangbiao651
// @param	ctxProvider		contexAPI.ClientProvider	客户端上下文
// @param	expectedPeers	int							组织的节点个数
// @return	peers			[]fabAPI.Peer				组织中的节点信息
// @return	err				error						获取节点过程中的错误信息
func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int, peerUrls map[string]bool) ([]fabAPI.Peer, error) {
	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		return nil, fmt.Errorf("error creating local context: %v", err)
	}

	// 获取组织中的节点
	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, fmt.Errorf("getting peers for MSP [%s] error: %v", ctx.Identifier().MSPID, serviceErr)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}
	ps := discoveredPeers.([]fabAPI.Peer)
	peers := []fabAPI.Peer{}
	for i, peer := range ps {
		if peerUrls[peer.URL()] {
			peers = append(peers, ps[i])
		}
	}

	return peers, nil
}

// @title InitService
// @description	 通过链码 ID ，通道 ID 组织信息以及 sdk	对通道信息进行初始化
func (t *SdkEnvInfo) InitService(chaincodeID, channelID string, org *OrgInfo, sdk *fabsdk.FabricSDK) error {
	handler := &SdkEnvInfo{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	var err error
	t.ChClient, err = channel.New(clientChannelContext)
	if err != nil {
		return err
	}
	t.EvClient, err = event.New(clientChannelContext, event.WithBlockEvents())
	if err != nil {
		return err
	}
	handler.ChClient = t.ChClient
	handler.EvClient = t.EvClient
	return nil
}

// @title regitserEvent
// @description 	注册链码事件
// @author	zhangbiao651
// @param	client		*event.Client		链码的管理客户端
// @param	chaincodeID	string				需要进行事件注册的链码ID
// @return	reg			fabAPI.Registration		链码事件的登记信息
// @return notifier		<-chan *fabAPI.CCEven	传递链码事件信息的通道
func regitserEvent(client *event.Client, chaincodeID string) (fabAPI.Registration, <-chan *fabAPI.CCEvent) {
	eventName := "chaincode-event"

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventName)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}

	return reg, notifier
}

// @title ChainCode EventListencer
// @description 对链码事件进行监听
// @author	zhangbiao651
// @param	c		*event.Client		链码事件管理客户端
// @param	ccID	string				链码所在的通道名
// @return 	reg		fabAPI.Registration	链码事件的注册信息
func ChainCodeEventListener(c *event.Client, ccID string) fabAPI.Registration {

	reg, notifier := regitserEvent(c, ccID)

	// consume event
	go func() {
		for e := range notifier {
			log.Printf("Receive cc event, ccid: %v \neventName: %v\n"+
				"payload: %v \ntxid: %v \nblock: %v \nsourceURL: %v\n",
				e.ChaincodeID, e.EventName, string(e.Payload), e.TxID, e.BlockNumber, e.SourceURL)
		}
	}()

	return reg
}

// @title	TxListener
// @description	对交易事件进行监听
// @author	zhangbiao651
// @param	c	*event.Clinet	链码事件管理客户端
// @param	txIDCh	chan string	传递交易信息的通道
func TxListener(c *event.Client, txIDCh chan string) {
	log.Println("Transaction listener start")
	defer log.Println("Transaction listener exit")

	for id := range txIDCh {
		// Register monitor transaction event
		log.Printf("Register transaction event for: %v", id)
		txReg, txCh, err := c.RegisterTxStatusEvent(id)
		if err != nil {
			log.Printf("Register transaction event error: %v", err)
			continue
		}
		defer c.Unregister(txReg)

		// Receive transaction event
		go func() {
			for e := range txCh {
				log.Printf("Receive transaction event: txid: %v, "+
					"validation code: %v, block number: %v",
					e.TxID,
					e.TxValidationCode,
					e.BlockNumber)
			}
		}()
	}
}

// @title BlockListener
// @description	对区块事件进行监听
// @param	ec	*event.Client	链码事件管理客户端
// @return	be	fabAPI.Registration	事件的登记信息
func BlockListener(ec *event.Client) fabAPI.Registration {
	// Register monitor block event
	beReg, beCh, err := ec.RegisterBlockEvent()
	if err != nil {
		log.Printf("Register block event error: %v", err)
	}
	log.Println("Registered block event")

	// Receive block event
	go func() {
		for e := range beCh {
			log.Printf("Receive block event:\nSourceURL: %v\nNumber: %v\nHash"+
				": %v\nPreviousHash: %v\n\n",
				e.SourceURL,
				e.Block.Header.Number,
				hex.EncodeToString(e.Block.Header.DataHash),
				hex.EncodeToString(e.Block.Header.PreviousHash))
		}
	}()

	return beReg
}

// @title getConfigFromChannel
// @description	同通道中获取配置文件
// @
func getConfigBlockFormChannel(channelID string, ordererEndpoint string, client *resmgmt.Client) (*cb.Config, error) {
	block, err := client.QueryConfigBlockFromOrderer(channelID, resmgmt.WithOrdererEndpoint(ordererEndpoint))
	if err != nil {
		return nil, err
	}

	blockDataEnvelope := &cb.Envelope{}
	err = proto.Unmarshal(block.Data.Data[0], blockDataEnvelope)
	if err != nil {
		return nil, err
	}

	blockDataPayload := &cb.Payload{}
	err = proto.Unmarshal(blockDataEnvelope.Payload, blockDataPayload)
	if err != nil {
		return nil, err
	}
	config := &cb.ConfigEnvelope{}
	err = proto.Unmarshal(blockDataPayload.Data, config)
	if err != nil {
		return nil, err
	}

	return config.Config, nil
}

func getSigningIdentity(sigIDPath string) *configtx.SigningIdentity {
	// Read certificate, private key and MSP ID from sigIDPath
	var (
		certificate *x509.Certificate
		privKey     crypto.PrivateKey
		mspID       string
		err         error
	)
	mspUser := filepath.Base(sigIDPath)

	certificate, err = readCertificate(filepath.Join(sigIDPath, "msp", "signcerts", fmt.Sprintf("%s-cert.pem", mspUser)))
	if err != nil {
		panic(err)
	}
	privKey, err = readPrivKey(filepath.Join(sigIDPath, "msp", "keystore", "priv_sk"))
	if err != nil {
		panic(err)
	}
	mspID = strings.Split(mspUser, "@")[1]
	return &configtx.SigningIdentity{
		Certificate: certificate,
		PrivateKey:  privKey,
		MSPID:       mspID,
	}
}

func readCertificate(certPath string) (*x509.Certificate, error) {
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(certBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("no PEM data found in cert[% x]", certBytes)
	}
	return x509.ParseCertificate(pemBlock.Bytes)
}

func readPrivKey(keyPath string) (crypto.PrivateKey, error) {
	privKeyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(privKeyBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("no PEM data found in private key[% x]", privKeyBytes)
	}
	return x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
}
