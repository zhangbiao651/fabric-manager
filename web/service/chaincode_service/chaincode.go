package chaincode_service

import "github.com/zhangbiao651/fabric-manager/web/models"

type ChaincodeService struct {
	ChaincodeID      string
	ChaincodePath    string
	ChaincodeVersion string

	PageNum  int
	PageSize int
}

func (service *ChaincodeService) GetAll() ([]*models.ChaincodeInfo, error) {
	chaincodes, err := models.GetChaincodes(service.PageNum, service.PageSize)
	if err != nil {
		return nil, err
	}
	return chaincodes, nil
}

func (service *ChaincodeService) Add() error {
	data := map[string]interface{}{
		"chaincode_id":      service.ChaincodeID,
		"chaincode_path":    service.ChaincodePath,
		"chaincode_version": service.ChaincodeVersion,
	}
	err := models.AddChaincode(data)
	if err != nil {
		return err
	}
	return nil
}

func (service *ChaincodeService) Get() (*models.ChaincodeInfo, error) {
	chaincode, err := models.GetChaincode(service.ChaincodeID)
	if err != nil {
		return nil, err
	}
	return chaincode, nil
}
