package models

import "github.com/jinzhu/gorm"

type OrdererInfo struct {
	gorm.Model
	OrdererAdminUser string // 排序组织的管理员用户
	OrdererOrgName   string // 排序组织的名称
	OrdererEndpoint  string // 排序组织的锚节点
	OrdererMspPath   string
}

func ExistOrdererByID(id int) (bool, error) {
	var orderer OrdererInfo
	err := db.Select("orderer_id").Where("id = ? ", id).First(&orderer).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if orderer.ID > 0 {
		return true, nil
	}
	return false, nil
}

// @title ExistOrdererByName
// @description	查看 Orderer 是否存在
// @author	zhangbiao651
// @param	id		int		待查询的排序组织ID
// @return  _		bool	排序组织的状态	true 表示已经存在 false 表示不存在
func ExistOrdererByName(name string) (bool, error) {
	var orderer OrdererInfo
	err := db.Select("id").Where("orderer_org_name = ?", name).First(&orderer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if orderer.ID > 0 {
		return true, nil
	}
	return false, nil
}

// @title GetOrderers
// @description 获得现有的所有排序组织信息
// @author	zhangbiao651
// @param	pageNum		int		页码
// @param	pageSize	int		页长
// @param   maps		interface{}
// @return	Orderers	[]*Orderer	所有的排序组织信息
// @return	err			error	获取排序组织信息过程中的错误信息
func GetOrderers(pageNum int, pageSize int) ([]*OrdererInfo, error) {
	var (
		orderers []*OrdererInfo
		err      error
	)
	if pageSize > 0 && pageNum > 0 {

		err = db.Find(&orderers).Offset(pageNum).Limit(pageSize).Error
	} else {

		err = db.Find(&orderers).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return orderers, nil

}

// @title GetOrdererTotal
// @descripion	获取排序组织的总数  用于分页时计算页码数量
// @author	zhangbiao651
// @param	maps	interface{} 	排序组织信息的映射
// @param	count	int				排序组织信息的总数
// @param	err		error			查询过程中的错误信息
func GetOrdererTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&OrdererInfo{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// @title GetOrdererById
// @description	通过排序组织的 ID 获取排序组织信息
// @author	zhangbiao651
// @param	id		int		需要查询的排序组织的 ID
// @return	Orderer	*Orderer	获取到的排序组织信息
func GetOrderer(name string) (*OrdererInfo, error) {
	var orderer OrdererInfo
	err := db.Where("orderer_org_name = ?  ", name).First(&orderer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &orderer, nil
}

// @title EditOrderer
// @description	更新 Orderer 的信息
// @author 	zhangbiao651
// @param	id		int		排序组织的 id 用于数据库中查询排序组织
// @param	data	interface 更新后的信息
// @return	err		error	更新过程中的错误信息
func EditOrderer(id int, data interface{}) error {
	if err := db.Model(&OrdererInfo{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// @title AddOrderer
// @description	向数据库中添加排序组织信息
// @author	zhangbiao651
// @param	data	map[string]interface{}	需要添加的排序组织信息
// @return	err		error		添加过程中发生的错误
func AddOrderer(data map[string]interface{}) error {
	orderer := OrdererInfo{
		OrdererOrgName:   data["orderer_org_name"].(string),
		OrdererEndpoint:  data["orderer_endpoint"].(string),
		OrdererMspPath:   data["orderer_msp_path"].(string),
		OrdererAdminUser: data["orderer_admin"].(string),
	}

	if err := db.Create(&orderer).Error; err != nil {
		return err
	}
	return nil
}

// @title DeleteOrderer
// @description	删除排序组织
// @author	zhangbiao651
// @param	id	int		需要删除的排序组织的 id
// @return	err	error	删除过程中发生的错误
func DeleteOrderer(id int) error {
	if err := db.Where("id = ?", id).Delete(&OrdererInfo{}).Error; err != nil {
		return err
	}
	return nil
}
