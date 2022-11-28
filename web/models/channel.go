package models

import (
	"fmt"

	"github.com/cloudflare/cfssl/log"
	"github.com/jinzhu/gorm"
)

type ChannelInfo struct {
	gorm.Model
	ChannelID string `gorm:"column:channel_id;size:20;not null;index:idx_ci" json:"channel_id"`
	// 通道名
	ChannelConfigPath string `gorm:"column:channel_config_path;size:255;not null" json:"channel_config_path"`

	// 通道的配置文件
	Org []*OrgInfo `gorm:"many2many:channel_orginfo"` // 通道内组织信息

}

func ExistChannelByName(name string) (bool, error) {
	var channel ChannelInfo
	err := db.Select("id").Where("channel_id = ?", name).First(&channel).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if channel.ID > 0 {
		return true, nil
	}
	return false, nil
}

// @title GetChannels
// @description 获得现有的所有排序组织信息
// @author	zhangbiao651
// @param	pageNum		int		页码
// @param	pageSize	int		页长
// @param   maps		interface{}
// @return	Channels	[]*Channel	所有的排序组织信息
// @return	err			error	获取排序组织信息过程中的错误信息
func GetChannels(pageNum int, pageSize int) ([]*ChannelInfo, error) {
	var (
		channels []*ChannelInfo
		err      error
	)
	if pageSize > 0 && pageNum > 0 {

		err = db.Model(&ChannelInfo{}).Preload("Org.Peers").Preload("Org").Find(&channels).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Model(&ChannelInfo{}).Preload("Org.Peers").Preload("Org").Find(&channels).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return channels, nil

}

// @title GetChannelTotal
// @descripion	获取排序组织的总数  用于分页时计算页码数量
// @author	zhangbiao651
// @param	maps	interface{} 	排序组织信息的映射
// @param	count	int				排序组织信息的总数
// @param	err		error			查询过程中的错误信息
func GetChannelTotal() (int, error) {
	var count int
	if err := db.Model(&ChannelInfo{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// @title GetChannelById
// @description	通过排序组织的 ID 获取排序组织信息
// @author	zhangbiao651
// @param	id		int		需要查询的排序组织的 ID
// @return	Channel	*Channel	获取到的排序组织信息
func GetChannel(name string) (*ChannelInfo, error) {
	var channel ChannelInfo
	err := db.Model(&ChannelInfo{}).Preload("Org").Preload("Org.Peers").First(&channel, "channel_id = ?", name).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	for i, org := range channel.Org {
		channel.Org[i].OrgAnchorFile = fmt.Sprintf("%s%s%sanchors.tx", org.OrgAnchorFile, channel.ChannelID, org.OrgName)
	}
	/*
		var orgs []*OrgInfo
		db.Model(&channel).Association("Org").Find(&orgs)
		for i := 0; i < len(orgs); i++ {

			org, err := GetOrgInfoByName(orgs[i].OrgName)
			if err != nil {
				return nil, err
			}
			orgs[i] = org
		}
		channel.Org = orgs
	*/
	return &channel, nil

}

// @title EditChannel
// @description	更新 Channel 的信息
// @author 	zhangbiao651
// @param	id		int		排序组织的 id 用于数据库中查询排序组织
// @param	data	interface 更新后的信息
// @return	err		error	更新过程中的错误信息
func EditChannel(id int, data interface{}) error {
	if err := db.Model(&ChannelInfo{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// @title AddChannel
// @description	向数据库中添加排序组织信息
// @author	zhangbiao651
// @param	data	map[string]interface{}	需要添加的排序组织信息
// @return	err		error		添加过程中发生的错误
func AddChannel(data map[string]interface{}) error {
	channel := ChannelInfo{
		ChannelID:         data["channel_id"].(string),
		ChannelConfigPath: data["channel_config_path"].(string),
	}

	if err := db.Create(&channel).Error; err != nil {
		return err
	}
	return nil
}

func OrgJoin(channelID string, orgnames []string) error {
	var channel ChannelInfo
	err := db.Where("channel_id = ?", channelID).First(&channel).Error
	if err != nil {
		return err
	}
	var orgs []*OrgInfo
	for _, ogn := range orgnames {
		log.Info("[models/channel.go/OrgJoin] 获取到的组织名为", ogn)
		org, err := GetOrgInfoByName(ogn)
		log.Info("[models/channel.go/OrgJoin]  获取到的组织信息为", org)
		if err != nil {
			return err
		}
		orgs = append(orgs, org)
	}
	channel.Org = orgs
	log.Info("[info models.Channel.go] 开始更新关联")
	err = db.Model(&channel).Association("Org").Append(orgs).Error
	if err != nil {
		log.Error("[ERROR  models.Channel.go],更新关联失败", err)
	}
	return nil
}

func OrgLeave(channelID string, orgnames []string) error {
	var channel ChannelInfo
	err := db.Where("channel_id = ?", channelID).First(&channel).Error
	if err != nil {
		return err
	}
	var orgs []*OrgInfo
	for _, ogn := range orgnames {
		log.Info("[models/channel.go/OrgJoin] 获取到的组织名为", ogn)
		org, err := GetOrgInfoByName(ogn)
		log.Info("[models/channel.go/OrgJoin]  获取到的组织信息为", org)
		if err != nil {
			return err
		}
		orgs = append(orgs, org)
	}
	channel.Org = orgs
	log.Info("[info models.Channel.go] 开始更新关联")
	err = db.Model(&channel).Association("Org").Delete(orgs).Error
	if err != nil {
		log.Error("[ERROR  models.Channel.go],更新关联失败", err)
	}
	return nil
}

// @title DeleteChannel
// @description	删除排序组织
// @author	zhangbiao651
// @param	id	int		需要删除的排序组织的 id
// @return	err	error	删除过程中发生的错误
func DeleteChannel(id int) error {
	if err := db.Where("id = ?", id).Delete(&ChannelInfo{}).Error; err != nil {
		return err
	}
	return nil
}
