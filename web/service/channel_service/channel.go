package channel_service

import "github.com/zhangbiao651/fabric-manager/web/models"

type ChannelService struct {
	ChannelID         string
	ChannelConfigPath string

	OrgNames []string

	PageNum  int
	PageSize int
}

func (service *ChannelService) Exist() (bool, error) {
	isExist, err := models.ExistChannelByName(service.ChannelID)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

func (service *ChannelService) Add() error {
	data := map[string]interface{}{
		"channel_id":          service.ChannelID,
		"channel_config_path": service.ChannelConfigPath,
	}
	if err := models.AddChannel(data); err != nil {
		return err
	}
	return nil
}

func (service *ChannelService) GetAll() ([]*models.ChannelInfo, error) {
	data, err := models.GetChannels(service.PageNum, service.PageSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (service *ChannelService) Get() (*models.ChannelInfo, error) {
	data, err := models.GetChannel(service.ChannelID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (service *ChannelService) AddOrg() error {
	if err := models.OrgJoin(service.ChannelID, service.OrgNames); err != nil {
		return err
	}
	return nil
}

func (service *ChannelService) MoveOrg() error {
	if err := models.OrgLeave(service.ChannelID, service.OrgNames); err != nil {
		return err
	}
	return nil
}
