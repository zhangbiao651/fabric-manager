package org_service

import (
	"github.com/cloudflare/cfssl/log"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/zhangbiao651/fabric-manager/web/models"
)

// 对 OrgService 模板进行封装 并添加了 用于分页的页数和页长
type OrgService struct {
	OrgAdminUser  string
	OrgName       string
	OrgMspId      string
	OrgUser       string
	OrgPeerNum    int
	OrgAnchorFile string

	OrgMspClient          *mspclient.Client          //组织的 MSP 客户端
	OrgAdminClientContext *contextAPI.ClientProvider // 组织管理员客户端的上下文
	OrgResMgmt            *resmgmt.Client
	Peers                 []*fab.Peer

	PageNum  int
	PageSize int
}

func (OrgService *OrgService) Exist() (bool, error) {
	isExit, err := models.ExistOrgInfoByOrgName(OrgService.OrgName)
	if err != nil {
		return false, err
	}
	log.Info("[orgService]组织的当前状态为:", isExit)
	log.Info("错误信息为", err)
	return isExit, nil

}

func (orgService *OrgService) Add() error {
	log.Info("[info] OrgService 正在添加组织 %s", orgService.OrgName)
	org := map[string]interface{}{
		"org_admin":       orgService.OrgAdminUser,
		"org_name":        orgService.OrgName,
		"org_msp_id":      orgService.OrgMspId,
		"org_user":        orgService.OrgUser,
		"org_peer_num":    orgService.OrgPeerNum,
		"org_anchor_file": orgService.OrgAnchorFile,
	}
	log.Info("[service/org.go] peernum = ", orgService.OrgPeerNum)
	isExit, err := models.ExistOrgInfoByOrgName(orgService.OrgName)
	if err != nil {
		return err
	}
	if isExit {
		err = models.EditOrgInfo(orgService.OrgName, org)
		if err != nil {
			return err
		}
		return nil
	}
	err = models.AddOrgInfo(org)
	if err != nil {
		return err
	}
	return nil
}

func (orgService *OrgService) Get() (*models.OrgInfo, error) {
	org, err := models.GetOrgInfoByName(orgService.OrgName)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return org, nil
}

func (orgService *OrgService) GetAll() ([]*models.OrgInfo, error) {

	orgs, err := models.GetOrgInfos(orgService.PageNum, orgService.PageSize)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (orgService *OrgService) Edit() error {
	org := map[string]interface{}{
		"org_admin":       orgService.OrgAdminUser,
		"org_name":        orgService.OrgName,
		"org_msp_id":      orgService.OrgMspId,
		"org_user":        orgService.OrgUser,
		"org_peer_num":    orgService.OrgPeerNum,
		"org_anchor_file": orgService.OrgAnchorFile,
	}

	err := models.EditOrgInfo(orgService.OrgName, org)
	if err != nil {
		return err
	}
	return nil
}
