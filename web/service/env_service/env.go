package env_service

import (
	"fmt"

	"github.com/cloudflare/cfssl/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/zhangbiao651/fabric-manager/sdkInit"
	"github.com/zhangbiao651/fabric-manager/web/service/chaincode_service"
	"github.com/zhangbiao651/fabric-manager/web/service/channel_service"
	"github.com/zhangbiao651/fabric-manager/web/service/orderer_service"
	"github.com/zhangbiao651/fabric-manager/web/service/org_service"
	"github.com/zhangbiao651/fabric-manager/web/service/peer_service"
)

type EnvService struct {
	ChannelID   string
	PeerUrl     []string
	ChaincodeID string
	ConfigPath  string
	OrdererName string

	Orgs        map[string]*sdkInit.OrgInfo
	SdkInfo     *sdkInit.SdkEnvInfo
	Sdk         *fabsdk.FabricSDK
	Application sdkInit.Application

	ChannelService   channel_service.ChannelService
	OrdererService   orderer_service.OrdererService
	OrgService       org_service.OrgService
	PeerService      peer_service.PeerService
	ChaincodeService chaincode_service.ChaincodeService
}

func (envService *EnvService) CreateChannel() error {
	peers := make(map[string]bool)
	for _, peer := range envService.PeerUrl {
		peers[peer] = true
	}
	log.Info("[service/env_service/env.go/CreateAndJoinChannel] sdkinfo = ", envService.SdkInfo, "    *sdkinfo= ", &envService.SdkInfo)
	err := sdkInit.CreateChannel(envService.SdkInfo)
	if err != nil {
		log.Error("[service/env_service/env.go] 获取通道失败")
		return fmt.Errorf("EnvService 创建通道失败 %v", err)
	}

	log.Info("[service/env_service/env.go] 获取通道成功")
	return nil
}

func (envService *EnvService) JoinChannel() error {
	peers := make(map[string]bool)
	for _, peer := range envService.PeerUrl {
		peers[peer] = true
	}
	err := sdkInit.JoinChannel(envService.SdkInfo, peers)
	if err != nil {
		return err
	}
	return nil
}

func (envService *EnvService) CreateCCLifecycle() error {
	peers := make(map[string]bool)
	for _, peer := range envService.PeerUrl {
		peers[peer] = true
	}

	err := sdkInit.CreateCCLifecycle(envService.SdkInfo, 1, envService.Sdk, peers)
	if err != nil {
		return fmt.Errorf("EnvService 创建链码生命周期失败：%v", err)
	}
	return nil
}

func (envService *EnvService) SetApplication() (*sdkInit.Application, error) {
	err := envService.SdkInfo.InitService(envService.ChaincodeID, envService.ChannelID, envService.SdkInfo.Orgs[envService.SdkInfo.OrgNames[0]], envService.Sdk)
	if err != nil {
		return nil, fmt.Errorf("EnvService 初始化 sdkInfoApplicaion 失败： %v", err)
	}
	return &envService.Application, nil
}

func (envService *EnvService) Init() error {

	// 获取通道信息
	envService.ChannelService = channel_service.ChannelService{
		ChannelID: envService.ChannelID,
	}
	channel, err := envService.ChannelService.Get()
	if err != nil {
		return fmt.Errorf("EnvService 获取通道信息失败：%v", err)
	}
	// 获取组织信息
	var ogns []string
	orgs := make(map[string]*sdkInit.OrgInfo)
	for _, org := range channel.Org {
		ogns = append(ogns, org.OrgName)
		orgInfo := &sdkInit.OrgInfo{
			OrgAdminUser:  org.OrgAdminUser,
			OrgName:       org.OrgName,
			OrgMspId:      org.OrgMspId,
			OrgUser:       org.OrgUser,
			OrgPeerNum:    org.OrgPeerNum,
			OrgAnchorFile: org.OrgAnchorFile,
		}
		orgs[org.OrgName] = orgInfo
	}

	// 获取排序节点信息
	envService.OrdererService = orderer_service.OrdererService{
		OrdererOrgName: envService.OrdererName,
	}
	ordererInfo, err := envService.OrdererService.Get()
	if err != nil {
		return fmt.Errorf("EnvService 获取排序节点信息失败")
	}

	// 获取链码信息
	envService.ChaincodeService = chaincode_service.ChaincodeService{
		ChaincodeID: envService.ChaincodeID,
	}
	chaincode, err := envService.ChaincodeService.Get()
	if err != nil {
		return fmt.Errorf("EnvService 获取链码信息失败")
	}

	// 创建交易信息
	sdkInfo := sdkInit.SdkEnvInfo{
		ChannelID:        channel.ChannelID,
		ChannelConfig:    channel.ChannelConfigPath,
		OrgNames:         ogns,
		Orgs:             orgs,
		OrdererAdminUser: ordererInfo.OrdererAdminUser,
		OrdererOrgName:   ordererInfo.OrdererOrgName,
		OrdererEndpoint:  ordererInfo.OrdererEndpoint,
		ChaincodeID:      chaincode.ChaincodeID,
		ChaincodePath:    chaincode.ChaincodePath,
		ChaincodeVersion: chaincode.ChaincodeVersion,
	}
	envService.SdkInfo = &sdkInfo

	log.Info("[env_service/env.go/init] config.yaml 的路径为", envService.ConfigPath)

	envService.Sdk, err = sdkInit.SetupSDK(envService.ConfigPath, &sdkInfo)

	if err != nil {
		log.Error("[env_service/env.go/Inif] ", err)
		return fmt.Errorf("EnvService 获取SDK 失败")
	}
	log.Info("[env_service/env.go/init] 环境初始化成功", sdkInfo)

	return nil
}
