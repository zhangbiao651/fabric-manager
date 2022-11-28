package sdkInit

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

func CreateChannel(info *SdkEnvInfo) error {
	log.Print("[INFO][sdkInit/sdkChennel.go/CreateChannel] 开始创建通道......")
	if len(info.Orgs) == 0 {
		log.Print("[ERROR][sdkInit/sdkChannel.go/CreateChannel] 创建通道失败:通道组织不能为空，请提供组织信息")
		return fmt.Errorf("通道组织不能为空，请提供组织信息")
	}

	// 获得所有组织的签名信息
	signIds := []msp.SigningIdentity{}
	log.Printf("[INFO][sdkInit/sdkChannel.go/CreateChannel] 获取组织签名信息中")
	for _, ogn := range info.OrgNames {
		org := info.Orgs[ogn]
		// Get signing identity that is used to sign create channel request
		orgSignId, err := org.OrgMspClient.GetSigningIdentity(org.OrgAdminUser)
		if err != nil {
			log.Printf("[ERROR][sdkInit/sdkChannel.go/CreateChannel] 通道创建失败：获取组织签名信息失败：%v", err)
			return fmt.Errorf("GetSigningIdentity error: %v", err)
		}
		signIds = append(signIds, orgSignId)
	}

	log.Printf("[INFO][sdkInit/sdkChannel.go/CreateChannel] 获取到组织签名信息,开始创建通道")
	// 创建通道
	if err := createChannel(signIds, info); err != nil {
		log.Printf("[ERROR][sdkInit/sdkChannel.go/CreateChannel] 创建通道失败，错误信息为：%v", err)
		return fmt.Errorf("Create channel error: %v", err)
	}

	log.Printf("[INFO][sdkInit/sdkChannel.go/CreateChannel] 创建通道成功")
	return nil

}

func JoinChannel(info *SdkEnvInfo, peerUrls map[string]bool) error {
	log.Print("[INFO][sdkInit/sdkChannel.go/JoinChannel]节点加入通道")
	for _, ogn := range info.OrgNames {
		org := info.Orgs[ogn]

		log.Print("[INFO][sdkInit/sdkChannel.go/JoinChannel]获取组织", org.OrgName, "下的节点")
		peers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, org.OrgPeerNum, peerUrls)
		if err != nil {
			log.Print("[ERROR][sdkInit/sdkChannel.go/JoinChannel]获取组织", org.OrgName, "下的节点失败")
			return err
		}
		log.Print("[INFO][sdkInit/sdkChannel.go/JoinChannel]获取组织", org.OrgName, "下的节点成功,其节点列表为", peers)
		for i, _ := range peers {
			log.Print("[INFO][sdkInit/sdkChannel.go/JoinChannel] 当前节点为", peers[i])
			if err := org.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(peers[i])); err != nil {
				log.Print("[ERROR][sdkInit/sdkChannel.go/JoinChannel] 节点", peers[i], "加入通道时发生错误：", err)
				return err
			}
			log.Print("节点[", peers[i].URL(), "]已经加入到通道当中")

		}

	}
	log.Print("[INFO][sdkInit/sdkChannel.go/JoinChannel] 所有节点已经加入通道")
	return nil
}

// @title createChannel
// @description  通过组织的签名以及通道信息完成通道的创建
// @author	zhangbiao651
// @param
func createChannel(signIDs []msp.SigningIdentity, info *SdkEnvInfo) error {
	// 通过 Channel 客户端完成通道的创建工作
	log.Printf("[INFO][sdkInit/sdkChennel/createChannel] 获取 resmgmt 客户端")
	log.Printf("[INFO][sdkInit/sdkChennel/createChannel] info.OrdererClientContext = %v", *info.OrdererClientContext)
	chMgmtClient, err := resmgmt.New(*info.OrdererClientContext)
	if err != nil {
		log.Print("[ERROR][sdkInit/sdkChennel/createChannel] 获取 resmgmt 客户端失败: ", err)
		return err
	}

	// 通过 通道的配置文件创建通道
	log.Print("[INFO][sdkInit/sdkChannel.go/createChannel] 读取配置文件,生成请求体")
	req := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: signIDs,
	}
	log.Printf("[INFO][sdkInit/sdkChannel.go/createChannel] 已生成请求体")

	//if _, err := chMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint)); err != nil {
	if _, err := chMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint)); err != nil {
		log.Printf("[ERROR][sdkInit/sdkChannel.go/createChannel] 创建通道时发生错误：%s", err.Error())
		return err
	}
	log.Printf("[INFO][sdkInit/sdkChannel.go/createdChannel] 使用管理员身份更新锚节点配置")
	for i, ogn := range info.OrgNames {
		org := info.Orgs[ogn]
		log.Print("[INFO][sdkInit/sdkChannel.go/createdChannel] 当前为", org, "创建更新锚节点请求结构体")
		req = resmgmt.SaveChannelRequest{
			ChannelID:         info.ChannelID,
			ChannelConfigPath: org.OrgAnchorFile,
			SigningIdentities: []msp.SigningIdentity{signIDs[i]},
		}
		log.Print("[INFO][sdkInit/sdkChannel.go/createChannel] 完成请求结构体创建，开始更新通道")
		if _, err = org.OrgResMgmt.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererEndpoint)); err != nil {
			log.Print("[error][sdkInit/sdkChannel.go/createChannel] 更新组织", ogn, "锚节点加入组织时发生错误，错误信息为", err)
			return fmt.Errorf("SaveChannel for anchor org %s error: %v", org.OrgName, err)
		}
	}
	log.Print("[INFO][sdkInit/sdkChannel.go/createChannel]使用每个org的管理员身份更新锚节点配置完成")

	return nil
}

/*
// @title OrgJoinChannel
// @description	向已经存在的通道加入组织
// @author	zhangbiao651
// @param	info		SdkEnvInfo		需要添加组织的通道的通道信息
// @param	org			*OrgInfo		需要向通道添加的组织信息
// @return	err			error			添加过程中的错误信息
func OrgJoinChannel(info SdkEnvInfo, orgInfo OrgInfo) error {
	orgConfig, err := ioutil.ReadFile(orgInfo.OrgConfigPath)
	if err != nil {
		log.Error(err)
	}
	var org configtx.Organization
	err = yaml.Unmarshal(orgConfig, org)
	baseConfig, err := getConfigBlockFormChannel(info.ChannelID, info.OrdererInfo.OrdererResMgmt)
	if err != nil {
		return fmt.Errorf("获取通道配置块错误，%v", err)
	}
	configTx := configtx.New(baseConfig)
	app := configTx.Application()

	err = app.SetOrganization(org)
	if err != nil {
		return err
	}
	configUpdateBytes, err := configTx.ComputeMarshaledUpdate(info.ChannelID)
	signingIdentity := getSigningIdentity(info.OrdererInfo.OrdererMspPath)
	configSignature, err := signingIdentity.CreateConfigSignature(configUpdateBytes)
	env, err := configtx.NewEnvelope(configUpdateBytes, configSignature)

	// Sign envelope
	envelopeSigningIdentity := getSigningIdentity(info.OrdererInfo.OrdererMspPath)
	err = envelopeSigningIdentity.SignEnvelope(env)
	envelopeBytes, err := proto.Marshal(env)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(info.ChannelConfigPath, envelopeBytes, 0640)

	r, err := os.Open(info.ChannelConfigPath)
	defer r.Close()
	resp, err := info.OrdererInfo.OrdererResMgmt.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     info.ChannelID,
		ChannelConfig: r,
	})
	if err != nil {
		return fmt.Errorf("更新通道失败： %v", err)
	}

	if resp.TransactionID == "" {
		return fmt.Errorf("更新通道失败:TransactionID 为空")
	}
	log.Infof("成功将组织 %s 加入通道", orgInfo.OrgName)

	return nil
}

// @title OrgLeaveChannel
// @decsription	删除通道中的组织
// @author	zhangbiao651
// @param	SdkEnvInfo		SdkEnvInfo		需要添加组织的通道的通道信息
// @param	org				*OrgInfo		需要向通道添加的组织信息
// @return	err				error			添加过程中的错误信息
func OrgLeaveChannel(info SdkEnvInfo, orgInfo OrgInfo) error {
	// 获取通道的配置块 并在移除其中的组织信息
	orgConfig, err := ioutil.ReadFile(orgInfo.OrgConfigPath)
	if err != nil {
		log.Error(err)
	}
	var org configtx.Organization
	err = yaml.Unmarshal(orgConfig, org)
	baseConfig, err := getConfigBlockFormChannel(info.ChannelID, info.OrdererInfo.OrdererResMgmt)
	if err != nil {
		return fmt.Errorf("获取通道配置块错误，%v", err)
	}
	configTx := configtx.New(baseConfig)
	app := configTx.Application()

	app.RemoveOrganization(orgInfo.OrgName)

	configUpdateBytes, err := configTx.ComputeMarshaledUpdate(info.ChannelID)
	signingIdentity := getSigningIdentity(info.OrdererInfo.OrdererMspPath)
	configSignature, err := signingIdentity.CreateConfigSignature(configUpdateBytes)
	env, err := configtx.NewEnvelope(configUpdateBytes, configSignature)

	// Sign envelope
	envelopeSigningIdentity := getSigningIdentity(info.OrdererInfo.OrdererMspPath)
	err = envelopeSigningIdentity.SignEnvelope(env)
	envelopeBytes, err := proto.Marshal(env)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(info.ChannelConfigPath, envelopeBytes, 0640)

	r, err := os.Open(info.ChannelConfigPath)
	defer r.Close()
	resp, err := info.OrdererInfo.OrdererResMgmt.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     info.ChannelID,
		ChannelConfig: r,
	})
	if err != nil {
		return fmt.Errorf("更新通道失败： %v", err)
	}

	if resp.TransactionID == "" {
		return fmt.Errorf("更新通道失败:TransactionID 为空")
	}
	log.Infof("成功将组织 %s 移出通道", orgInfo.OrgName)

	return nil

}
*/
