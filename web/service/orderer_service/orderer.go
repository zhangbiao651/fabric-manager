package orderer_service

import "github.com/zhangbiao651/fabric-manager/web/models"

type OrdererService struct {
	ID               int
	OrdererAdminUser string // 排序组织的管理员用户
	OrdererOrgName   string // 排序组织的名称
	OrdererEndpoint  string // 排序组织的锚节点
	OrdererMspPath   string

	PageNum  int
	PageSize int
}

func (service *OrdererService) Add() error {
	data := map[string]interface{}{
		"orderer_admin":    service.OrdererAdminUser,
		"orderer_org_name": service.OrdererOrgName,
		"orderer_endpoint": service.OrdererEndpoint,
		"orderer_msp_path": service.OrdererMspPath,
	}
	err := models.AddOrderer(data)
	if err != nil {
		return err
	}

	return nil

}

func (service *OrdererService) Exist() (bool, error) {
	isExist, err := models.ExistOrdererByName(service.OrdererOrgName)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

func (service *OrdererService) GetAll() ([]*models.OrdererInfo, error) {
	orderers, err := models.GetOrderers(service.PageNum, service.PageSize)
	if err != nil {
		return nil, err
	}
	return orderers, nil
}

func (service *OrdererService) Get() (*models.OrdererInfo, error) {
	orderer, err := models.GetOrderer(service.OrdererOrgName)
	if err != nil {
		return nil, err
	}
	return orderer, nil
}
