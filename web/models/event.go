package models

import "gorm.io/gorm"

// 事件信息
type EventInfo struct {
	gorm.Model
	// 通道信息
	// 通道ID
	ChannelID string `gorm:"column:channel_id;size:20;not null;index:idx_ci" json:"channel_id"`
	// 通道配置文件的路径
	ChannelConfig string `gorm:"column:channel_config_path;size:255;not null" json:"channel_config_path"`
	Peers []string `json:""`
	// 排序服务节点信息
	// 排序节点的管理员用户
	OrdererAdminUser string `gorm:"column:org_admin;size:20;not null" json:"org_admin"`
	// 排序节点组织的组织名称
	OrdererOrgName string `gorm:"column:orderer_org_name;size:20;not null" json:"orderer_org_name"`
	// 排序节点组织的域名
	OrdererEndpoint string `gorm:"column:orderer_endpoint;size:100;not null" json:"orderer_endpoint"`

	// 链码信息
	// 链码 ID
	ChaincodeID string `gorm:"column:chaincode_id;size:20;" json:"chaincode_id"`
	// 链码的路径
	ChaincodePath    string `gorm:"column:chaincode_path;size:255;" json:"chaincode_path"`
	ChaincodeVersion string `gorm:"column:chaincode_version;size:20;l" json:"chaincode_version"`
}

// @title ExisEventInfo
// @descriotion  通过通道名称查看事件信息是否已经存在
// @author zhangbiao651
// @param		channelID		string		通道的 ID
// @return		_				bool		事件信息状态  true 表示以存在  false表示不存在
func ExistEventInfoByChannelID(channelID string) (bool, error) {
	var EventInfo EventInfo
	err := db.Select("channel_id").Where("channel_id = ? AND deleted_at = nll", channelID).First(&EventInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if EventInfo.ID > 0 {
		return true, nil
	}

	return false, nil
}

// @title GetChannelsInfo
// @description 分页查询组织信息
// @author zhangbiao651
// @param		pageNum			int			页面数
// @param		pageSize		int			页长
// @return		EventInfo		[]*EventInfo	该页所有组织的组织信息
func GetChannelsInfo(pageNum int, pageSize int) ([]*EventInfo, error) {
	var EventInfo []*EventInfo
	err := db.Where("delete_at = null").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&EventInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return EventInfo, nil
}

// @title  AddEventInfo
// @description  向数据库中添加组织信息
// @author zhangbiao651
// @param		EventInfo		EventInfo		需要记录的组织信息
// @return 		err			error		记录组织信息时发生的错误
func AddEventInfo(EventInfo EventInfo) error {
	err := db.Create(&EventInfo).Error
	if err != nil {
		return err
	}
	return nil
}

// @title UpdateEventInfo
// @description  将更新后的组织信息同步更新到数据库当中
// @author zhangbiao651
// @param		EventInfo		EventInfo		需要同步到数据库的组织信息
// @return		err			error		更新组织信息时发生的错误
func UpdateEventInfo(eventInfo EventInfo) error {
	if err := db.Model(&EventInfo{}).Where("id = ? ", eventInfo.ID, 0).Update(eventInfo).Error; err != nil {
		return err
	}

	return nil
}

// @title DeleteEventInfo
// @description  删除数据库中的组织信息
// @author zhangbiao651
// @param		EventInfo		EventInfo		需要删除的组织信息
// @return		err			error		删除数据过程中发生的错误
func DeleteEventInfo(eventInfo EventInfo) error {
	if err := db.Where("id = ?", eventInfo.ID).Delete(EventInfo{}).Error; err != nil {
		return err
	}
	return nil
}

// @title GetPageNum
// @description   获取该项数据的页码总数
// @author zhangbiao651
// @param		pageSize		int		数据的页容量
// @return 		pageNum		int		数据可以被分为多少页
func GetPageNums(pageSize int) (int, error) {
	var total int
	err := db.Model(&EventInfo{}).Count(&total).Error
	pageNum := total / pageSize
	if total%pageSize != 0 {
		pageNum++
	}

	return pageNum, err

}
