package sdkInit

import (
	"fmt"
	"log"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//signaturending approve transaction proposal failed to verify signature

// @title SetupSDK
// @description 通过配置文件初始化 Fabric SDK
// @author  	zhangbiao651
// @param		configFile		string    	config.yaml的路径
// @return		err				error		实例化 Fabric SDK 过程中的错误信息
func SetupSDK(configFile string, info *SdkEnvInfo) (*fabsdk.FabricSDK, error) {
	var err error
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Print("[sdkInit/asdSetup.go/SetupSDK] 实例化Fabric SDK 失败", err)
		return nil, fmt.Errorf("实例化Fabric SDK 失败： %v\n", err)
	}
	log.Printf("[sdkInit/asdSetup.go/SetupSDK]Fabric SDK初始化成功")

	// 为组织获得Client句柄和Context信息
	for _, org := range info.Orgs {
		org.OrgMspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(org.OrgName))
		if err != nil {
			return nil, err
		}
		orgContext := sdk.Context(fabsdk.WithUser(org.OrgAdminUser), fabsdk.WithOrg(org.OrgName))
		org.OrgAdminClientContext = &orgContext
		log.Print("[INFO][sdkInit/sdkSetup.go/SetupSDK]，组织[", org.OrgName, "]的用户为[", org.OrgAdminUser, "]")

		// New returns a resource management client instance.
		resMgmtClient, err := resmgmt.New(orgContext)
		if err != nil {
			return nil, fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
		}
		org.OrgResMgmt = resMgmtClient
		org.Peers = make(map[string]*fabAPI.Peer)
	}
	ordererClientContext := sdk.Context(fabsdk.WithUser(info.OrdererAdminUser), fabsdk.WithOrg(info.OrdererOrgName))
	info.OrdererClientContext = &ordererClientContext
	log.Printf("[sdkInit/sdkSetup.go/setup],ordererClientContext = [%v],info.OrdererClientContext = [%v]", &ordererClientContext, info.OrdererClientContext)

	return sdk, nil
}

/*
// @title UpDateSDK
// @description 与SetupSDK 逻辑相同 不同名称
// @author	zhangbiao651
// @param		configFile		string		config.yaml的路径
// @return		err				error		升级 Fabric SDK 过程中的错错误信息
func UpdateSDK(configFile string, info SdkEnvInfo) error {
	fabSdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		return fmt.Errorf("根据新配置文件生成新 FabricSDK 失败")
	}
	orderer := info.OrdererInfo
	// 通过排序节点为通道获取 Client 和 Context信息
	orderer.OrdererMspClient, err = mspclient.New(fabSdk.Context(), mspclient.WithOrg(orderer.OrdererOrgName))
	if err != nil {
		return err
	}
	orgContext := fabSdk.Context(fabsdk.WithUser(orderer.OrdererAdminUser), fabsdk.WithOrg(orderer.OrdererOrgName))
	orderer.OrdererClientContext = &orgContext
	resMgmtClient, err := resmgmt.New(orgContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端 Context 创建通道管理客户端失败：%v", err)
	}
	orderer.OrdererResMgmt = resMgmtClient
	ordererClientContext := fabSdk.Context(fabsdk.WithUser(orderer.OrdererAdminUser), fabsdk.WithOrg(orderer.OrdererOrgName))
	orderer.OrdererClientContext = &ordererClientContext
	info.OrdererInfo = orderer

	// 为组织获得Client句柄和Context信息
	for _, org := range info.Orgs {
		org.OrgMspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(org.OrgName))
		if err != nil {
			return err
		}
		orgContext := sdk.Context(fabsdk.WithUser(org.OrgAdminUser), fabsdk.WithOrg(org.OrgName))
		org.OrgAdminClientContext = &orgContext

		// New returns a resource management client instance.
		resMgmtClient, err := resmgmt.New(orgContext)
		if err != nil {
			return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
		}
		org.OrgResMgmt = resMgmtClient
	}
	sdk = fabSdk
	log.Info("Fabric SDK 升级成功")
	return nil
}
*/
