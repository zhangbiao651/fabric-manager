package sdkInit

import (
	"fmt"
	"strings"

	"github.com/cloudflare/cfssl/log"
	mb "github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/policydsl"
)

// @title CreateCCLifecycle
// @description	通过 FabricSDK 为通道安装链码
// @param		info		*ChannelInfo		通道信息
// @param		sequence	int64				序列码
// @param		sdk			*fabsdk.FabricSDK	FabricSDK实例
// @return		err			error				安装链码时的错误信息
func CreateCCLifecycle(info *SdkEnvInfo, sequence int64, sdk *fabsdk.FabricSDK, peerUrls map[string]bool) error {
	if len(info.Orgs) == 0 {
		return fmt.Errorf("组织数不应该为0")
	}
	// Package cc
	fmt.Println(">> 开始打包链码......")
	label, ccPkg, err := packageCC(info.ChaincodeID, info.ChaincodeVersion, info.ChaincodePath)
	if err != nil {
		return fmt.Errorf("pakcagecc error: %v", err)
	}
	packageID := lcpackager.ComputePackageID(label, ccPkg)
	fmt.Println(">>>>>>>>>>>>>>>>>PackageID = [", packageID, "]")
	fmt.Println(">> 打包链码成功")

	// Install cc
	fmt.Println(">> 开始安装链码......")
	if err := installCC(label, ccPkg, info.Orgs, peerUrls); err != nil {
		return fmt.Errorf("installCC error: %v", err)
	}

	// Get installed cc package
	if err := getInstalledCCPackage(packageID, info.Orgs[info.OrgNames[0]], peerUrls); err != nil {
		return fmt.Errorf("getInstalledCCPackage error: %v", err)
	}

	// Query installed cc
	if err := queryInstalled(packageID, info.Orgs[info.OrgNames[0]], peerUrls); err != nil {
		return fmt.Errorf("queryInstalled error: %v", err)
	}
	fmt.Println(">> 安装链码成功")

	// Approve cc
	fmt.Println(">> 组织认可智能合约定义......")
	if err := approveCC(packageID, info.ChaincodeID, info.ChaincodeVersion, sequence, info.ChannelID, info.OrgNames, info.Orgs, info.OrdererEndpoint, peerUrls); err != nil {
		return fmt.Errorf("approveCC error: %v", err)
	}

	fmt.Println(">> 查询认可结果")
	// Query approve cc
	if err := queryApprovedCC(info.ChaincodeID, sequence, info.ChannelID, info.OrgNames, info.Orgs, peerUrls); err != nil {
		return fmt.Errorf("queryApprovedCC error: %v", err)
	}
	fmt.Println(">> 组织认可智能合约定义完成")

	// Check commit readiness
	fmt.Println(">> 检查智能合约是否就绪......")
	if err := checkCCCommitReadiness(packageID, info.ChaincodeID, info.ChaincodeVersion, sequence, info.ChannelID, info.OrgNames, info.Orgs, peerUrls); err != nil {
		return fmt.Errorf("checkCCCommitReadiness error: %v", err)
	}
	fmt.Println(">> 智能合约已经就绪")

	// Commit cc
	fmt.Println(">> 提交智能合约定义......")
	if err := commitCC(info.ChaincodeID, info.ChaincodeVersion, sequence, info.ChannelID, info.OrgNames, info.Orgs, info.OrdererEndpoint); err != nil {
		return fmt.Errorf("commitCC error: %v", err)
	}
	// Query committed cc
	if err := queryCommittedCC(info.ChaincodeID, info.ChannelID, sequence, info.OrgNames, info.Orgs, peerUrls); err != nil {
		return fmt.Errorf("queryCommittedCC error: %v", err)
	}
	fmt.Println(">> 智能合约定义提交完成")

	// Init cc
	//	fmt.Println(">> 调用智能合约初始化方法......")
	//	if err := initCC(info.ChaincodeID, info.ChannelID, info.Orgs[info.OrgNames[0]], sdk); err != nil {
	//		return fmt.Errorf("initCC error: %v", err)
	//	}
	//	fmt.Println(">> 完成智能合约初始化")
	return nil
}

// @title	packageCC
// @description	打包链码
// @param	ccName		string	链码名称
// @param	ccVersion	string	链码的版本
// @param	ccpath		string	链码的存储路径
// @return	lable
func packageCC(ccName, ccVersion, ccpath string) (string, []byte, error) {
	label := ccName + "_" + ccVersion
	desc := &lcpackager.Descriptor{
		Path:  ccpath,
		Type:  pb.ChaincodeSpec_GOLANG,
		Label: label,
	}
	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		return "", nil, fmt.Errorf("Package chaincode source error: %v", err)
	}
	return desc.Label, ccPkg, nil
}

// @title installCC
// @description		将打包后的链码安装到节点中去
// @param	label	string		打包链码后生成的标签
// @param	ccPkg	[]byte		链码打包后的字节数组
// @param	orgs	map[string]*Orginfo		需要安装链码的节点
// @return	err		error		链码安装过程中的错误
func installCC(label string, ccPkg []byte, orgs map[string]*OrgInfo, peerUrls map[string]bool) error {

	//
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	packageID := lcpackager.ComputePackageID(installCCReq.Label, installCCReq.Package)
	for _, org := range orgs {
		log.Info("[sdkInit/sdkChaincode.go/installcc] 为组织[", org.OrgName, "]安装链码，查询节点")
		orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, org.OrgPeerNum, peerUrls)
		if err != nil {
			log.Error("[ERROR][sdkInit/sdkChaincode.go/installcc]组织[", org.OrgName, "]获取节点失败")
			return fmt.Errorf("DiscoverLocalPeers error: %v", err)
		}
		log.Info("[INFO][sdkInit/sdkChaincode.go/installcc]组织[", org.OrgName, "],下的节点为", orgPeers)
		for _, peer := range orgPeers {
			log.Info("[INFO][sdkInit/sdkChaincode.go/installcc]: 检查[", peer.URL(), "]是否安装链码")
			if flag, _ := checkInstalled(packageID, peer, org.OrgResMgmt); flag == false {
				log.Info("[INFO][sdkInit/sdkChaincode.go/installcc]: 检查到[", peer.URL(), "]尚未安装链码，开始安装")
				if _, err := org.OrgResMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithTargets(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
					log.Error()
					return fmt.Errorf("LifecycleInstallCC error: %v", err)
				}
			}
		}

	}
	return nil
}

// @title getInstalledCCPackage
// @description	检索给定包 ID 的已安装链代码包
// @author	zhangbiao651
// @param	packageID	string		链码的包 ID
// @param	org			*OrgInfo	检索的组织信息
// @return	err			error		检索过程中的错误信息
func getInstalledCCPackage(packageID string, org *OrgInfo, peerUrls map[string]bool) error {
	// use org1
	orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, 1, peerUrls)
	if err != nil {
		return fmt.Errorf("DiscoverLocalPeers error: %v", err)
	}

	if _, err := org.OrgResMgmt.LifecycleGetInstalledCCPackage(packageID, resmgmt.WithTargets([]fab.Peer{orgPeers[0]}...)); err != nil {
		return fmt.Errorf("LifecycleGetInstalledCCPackage error: %v", err)
	}
	return nil
}

// @title	queryInstalled
// @description	查询链码的安装状态
// @author	zhangbiao651
// @param	packageID	string	打包后的链码的包ID
// @param	org			*ORgInfo	组织信息
// @return	err			error	检查过程中的错误信息
func queryInstalled(packageID string, org *OrgInfo, peerUrls map[string]bool) error {
	orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, 1, peerUrls)
	if err != nil {
		return fmt.Errorf("DiscoverLocalPeers error: %v", err)
	}
	resp1, err := org.OrgResMgmt.LifecycleQueryInstalledCC(resmgmt.WithTargets([]fab.Peer{orgPeers[0]}...))
	if err != nil {
		return fmt.Errorf("LifecycleQueryInstalledCC error: %v", err)
	}
	packageID1 := ""
	for _, t := range resp1 {
		if t.PackageID == packageID {
			packageID1 = t.PackageID
		}
	}
	if !strings.EqualFold(packageID, packageID1) {
		return fmt.Errorf("check package id error")
	}
	return nil
}

// @title checkInstalled
// @description	对链码的安装结果进行检查
// @param	packageID	string		链码打包后的包 ID
// @param	peer		fab.Peer	通道中的节点
// @param	client		*resmgmt.Client	对通道和链码进行管理的客户端
// @return	flag		bool		链码的安装状态	true 表现以及安装 false	表示未安装
func checkInstalled(packageID string, peer fab.Peer, client *resmgmt.Client) (bool, error) {
	flag := false
	resp1, err := client.LifecycleQueryInstalledCC(resmgmt.WithTargets(peer))
	if err != nil {
		return flag, fmt.Errorf("LifecycleQueryInstalledCC error: %v", err)
	}
	for _, t := range resp1 {
		if t.PackageID == packageID {
			flag = true
		}
	}
	return flag, nil
}

// @title approveCC
// @description  批准链码的安装
// @author	zhangbiao651
// @param	packageID		string			生成批准请求的参数
// @param	sequence		int				生成批准请求的参数
// @param	chaincode		ChaincodeInfo	链码信息（主要使用到链码的 ID 以及版本）
// @param	channelID		string			通道的名称
// @param	orgs			map[string]*OrgInfo	通道中的组织
// @param	ordererEndpoint	string			排序节点的锚节点
func approveCC(packageID string, chaincodeID, chaincodeVersion string, sequence int64, channelID string, orgnames []string, orgs map[string]*OrgInfo, ordererEndpoint string, peerUrls map[string]bool) error {

	mspIDs := []string{}
	// 获取组织 MSPID
	for _, ogn := range orgnames {
		org := orgs[ogn]
		mspIDs = append(mspIDs, org.OrgMspId)
		log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC] 组织[", org.OrgName, "]的MSPID 为[", org.OrgMspId, "]")
	}
	log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC] 获取到组织的MSPID", mspIDs)

	// 策略获取
	log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC] 获取策略")
	ccPolicy := policydsl.SignedByNOutOfGivenRole(int32(len(mspIDs)), mb.MSPRole_MEMBER, mspIDs)
	log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC] 获取策略完成,策略如下,", ccPolicy.String())

	log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC] 生成链码批准结构体")
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              chaincodeID,
		Version:           chaincodeVersion,
		PackageID:         packageID,
		Sequence:          sequence,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC] 链码批准结构体成功生成", approveCCReq)
	for _, ogn := range orgnames {
		org := orgs[ogn]
		orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, org.OrgPeerNum, peerUrls)
		fmt.Printf(">>> chaincode approved by %s peers:\n", org.OrgName)

		for _, p := range orgPeers {
			fmt.Printf("	%s\n", p.URL())
		}

		if err != nil {
			return fmt.Errorf("DiscoverLocalPeers error: %v", err)
		}
		log.Info("[INFO][sdkInit/sdkChaincode.go/approveCC]  [", ogn, "] 开始链码批准")
		if _, err := org.OrgResMgmt.LifecycleApproveCC(channelID, approveCCReq, resmgmt.WithTargets(orgPeers...), resmgmt.WithOrdererEndpoint(ordererEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
			log.Error("[Error][sdkInit/sdkChaincode.go/approveCC]  [", ogn, "] 链码批准失败", err)
			return fmt.Errorf("LifecycleApproveCC error: %v", err)
		}

	}
	return nil
}

// @title	queryApprovedCC
// @description	查询组织是否认可智能合约
// @author	zhangbiao651
// @param	ccName		string	链码名称
// @param	sequence	int64	序列号
// @param	channelID	string	通道ID
// @param	orgs		map[string]*orgInfo	该通道中组织的组织信息
// @return  err			error	查询过程中的错误信息
func queryApprovedCC(ccName string, sequence int64, channelID string, orgnames []string, orgs map[string]*OrgInfo, peerUrls map[string]bool) error {
	// 创建查询结构体
	queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     ccName,
		Sequence: sequence,
	}

	// 查询所有节点是否提交智能合约
	log.Info("[INFO][sdkInit/sdkChaincode.go] 开始检查组织是否认可合约")
	for _, ogn := range orgnames {
		org := orgs[ogn]

		log.Info("[INFO][sdkInit/sdkChaincode.go] 当前查询的组织为", ogn)
		orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, org.OrgPeerNum, peerUrls)
		if err != nil {
			return fmt.Errorf("DiscoverLocalPeers error: %v", err)
		}
		// 查询组织中的每个节点的情况
		for _, p := range orgPeers {

			log.Info("[INFO][sdkInit/sdkChaincode.go] 当前查询的节点为", p.URL())
			resp, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
				func() (interface{}, error) {
					resp1, err := org.OrgResMgmt.LifecycleQueryApprovedCC(channelID, queryApprovedCCReq, resmgmt.WithTargets(p))
					if err != nil {
						return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleQueryApprovedCC returned error: %v", err), nil)
					}
					return resp1, err
				},
			)
			if err != nil {
				return fmt.Errorf("Org %s Peer %s NewInvoker error: %v", org.OrgName, p.URL(), err)
			}
			if resp == nil {
				return fmt.Errorf("Org %s Peer %s Got nil invoker", org.OrgName, p.URL())
			}
		}
	}
	return nil
}

// @title	checkCCCommitReadiness
// @description	检查智能合约是否就绪
// @author	zhangbiao651
// @param	packageID	string	打包后智能合约的 ID
// @param	ccName		string	链码的名称
// @param	ccVersion	string	链码的版本
// @param	sequence	int64	序列号
// @param	channelID	string	通道名称
// @param	orgs		msp[string]*OrgInfo		通道内所有组织的信息
// @return	err			error	检查过程中产生的错误信息
func checkCCCommitReadiness(packageID string, chaincodeID, chaincodeVersion string, sequence int64, channelID string, orgnames []string, orgs map[string]*OrgInfo, peerUrls map[string]bool) error {
	mspIds := []string{}
	for _, ogn := range orgnames {
		org := orgs[ogn]
		mspIds = append(mspIds, org.OrgMspId)
	}
	ccPolicy := policydsl.SignedByNOutOfGivenRole(int32(len(mspIds)), mb.MSPRole_MEMBER, mspIds)
	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:    chaincodeID,
		Version: chaincodeVersion,
		//PackageID:         packageID,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		Sequence:          sequence,
		InitRequired:      true,
	}

	for _, ogn := range orgnames {
		org := orgs[ogn]
		orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, org.OrgPeerNum, peerUrls)
		if err != nil {
			return fmt.Errorf("DiscoverLocalPeers error: %v", err)
		}
		for _, p := range orgPeers {
			resp, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
				func() (interface{}, error) {
					resp1, err := org.OrgResMgmt.LifecycleCheckCCCommitReadiness(channelID, req, resmgmt.WithTargets(p))
					fmt.Printf("LifecycleCheckCCCommitReadiness cc = %v, = %v\n", chaincodeID, resp1)
					if err != nil {
						return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleCheckCCCommitReadiness returned error: %v", err), nil)
					}
					flag := true
					for _, r := range resp1.Approvals {
						flag = flag && r
					}
					if !flag {
						return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleCheckCCCommitReadiness returned : %v", resp1), nil)
					}
					return resp1, err
				},
			)
			if err != nil {
				return fmt.Errorf("NewInvoker error: %v", err)
			}
			if resp == nil {
				return fmt.Errorf("Got nill invoker response")
			}
		}
	}

	return nil
}

// @title commitCC
// @description 向各个组织提交智能合约
// @param		info		ChaincodeInfo	智能合约的信息 主要使用到其 ID 以及版本
// @param		sequence	int64		序列玛
// @param		channelID	string		通道的名称
// @param		orgname		string		节点中的某个组织
// @param		orgs		map[string]*OrgInfo		通道中的组织列表
// @param		ordererEndpoint		string			排序组织的锚节点
// @return		err			error		提交过程中的错误信息
func commitCC(chaincodeID, chaincodeVersion string, sequence int64, channelID string, orgnames []string, orgs map[string]*OrgInfo, ordererEndpoint string) error {
	mspIDs := []string{}
	// 获取组织中的 MSPID
	for _, ogn := range orgnames {
		org := orgs[ogn]
		mspIDs = append(mspIDs, org.OrgMspId)
	}

	// 获取链码的策略
	ccPolicy := policydsl.SignedByNOutOfGivenRole(int32(len(mspIDs)), mb.MSPRole_MEMBER, mspIDs)

	// 创建链码提交发送体
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              chaincodeID,
		Version:           chaincodeVersion,
		Sequence:          sequence,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	var err error
	for _, ogn := range orgnames {
		org := orgs[ogn]
		_, err = org.OrgResMgmt.LifecycleCommitCC(channelID, req, resmgmt.WithOrdererEndpoint(ordererEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		break
	}
	if err != nil {
		return fmt.Errorf("LifecycleCommitCC error: %v", err)
	}
	return nil
}

// @title	queryCommittedCC
// @description 查询智能合约是否安装完成
// @author	zhangbiao651
// @param	ccName		string  链码名称
// @param	channelID	string	通道ID
// @param	sequence	int65	区块链高度
// @param	orgs		map[string]*OrgInfo	节点中的组织信息
// @return	err			error	查询过程中发生的错误
func queryCommittedCC(ccName string, channelID string, sequence int64, orgnames []string, orgs map[string]*OrgInfo, peerUrls map[string]bool) error {
	req := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: ccName,
	}

	// 通过组织管理元客户端以及节点数量查询所有节点
	for _, ogn := range orgnames {
		org := orgs[ogn]
		orgPeers, err := DiscoverLocalPeers(*org.OrgAdminClientContext, org.OrgPeerNum, peerUrls)
		if err != nil {
			return fmt.Errorf("DiscoverLocalPeers error: %v", err)
		}

		// 分别查询每个节点
		for _, p := range orgPeers {
			resp, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
				func() (interface{}, error) {
					resp1, err := org.OrgResMgmt.LifecycleQueryCommittedCC(channelID, req, resmgmt.WithTargets(p))
					if err != nil {
						return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleQueryCommittedCC returned error: %v", err), nil)
					}
					flag := false
					for _, r := range resp1 {
						if r.Name == ccName && r.Sequence == sequence {
							flag = true
							break
						}
					}
					if !flag {
						return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleQueryCommittedCC returned : %v", resp1), nil)
					}
					return resp1, err
				},
			)
			if err != nil {
				return fmt.Errorf("NewInvoker error: %v", err)
			}
			if resp == nil {
				return fmt.Errorf("Got nil invoker response")
			}
		}
	}
	return nil
}

// @title initcc
// @description  初始化链码
// @param	ccName		string		链码 ID
// @param	channelID	string		通道 ID
// @param	org			*OrgInfo	组织信息
// @param	sdk			*fabsdk.FabricSDK		FabricSDK 实例
// @return	err			error		初始化链码过程中的错误信息
func initCC(ccName string, channelID string, org *OrgInfo, sdk *fabsdk.FabricSDK) error {
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return fmt.Errorf("Failed to create new channel client: %s", err)
	}

	// init
	_, err = client.Execute(channel.Request{ChaincodeID: ccName, Fcn: "init", Args: nil, IsInit: true},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return fmt.Errorf("Failed to init: %s", err)
	}
	return nil
}
