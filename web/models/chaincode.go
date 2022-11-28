package models

import "github.com/jinzhu/gorm"

type ChaincodeInfo struct {
	gorm.Model
	ChaincodeID string `gorm:"column:chaincode_id;size:20;not null" json:"chaincode_id"`
	// 链码的路径
	ChaincodePath    string `gorm:"column:chaincode_path;size:255255;not null" json:"chaincode_path"`
	ChaincodeVersion string `gorm:"column:chaincode_version;size:20;not null" json:"chaincode_version"`
}

// @title ExistChaincodeByID
// @description	查看 Chaincode 是否存在
// @author	zhangbiao651
// @param	id		int		待查询的链码ID
// @return  _		bool	链码的状态	true 表示已经存在 false 表示不存在
func ExistChaincodeByID(id int) (bool, error) {
	var chaincode ChaincodeInfo
	err := db.Select("chaincode_id").Where("id = ? AND delete_at = nil", id).First(&chaincode).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if chaincode.ID > 0 {
		return true, nil
	}
	return false, nil
}

// @title ExistChaincodeByName
// @description	查看 Chaincode 是否存在
// @author	zhangbiao651
// @param	id		int		待查询的链码ID
// @return  _		bool	链码的状态	true 表示已经存在 false 表示不存在
func ExistChaincodeByName(name string) (bool, error) {
	var chaincode ChaincodeInfo
	err := db.Select("id").Where("chaincode_id = ? ", name).First(&chaincode).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if chaincode.ID > 0 {
		return true, nil
	}
	return false, nil
}

// @title GetChaincodes
// @description 获得现有的所有链码信息
// @author	zhangbiao651
// @param	pageNum		int		页码
// @param	pageSize	int		页长
// @param   maps		interface{}
// @return	chaincodes	[]*Chaincode	所有的链码信息
// @return	err			error	获取链码信息过程中的错误信息
func GetChaincodes(pageNum int, pageSize int) ([]*ChaincodeInfo, error) {
	var (
		chaincodes []*ChaincodeInfo
		err        error
	)
	if pageSize > 0 && pageNum > 0 {

		err = db.Find(&chaincodes).Offset(pageNum).Limit(pageSize).Error
	} else {

		err = db.Find(&chaincodes).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return chaincodes, nil
}

// @title GetChaincodeTotal
// @descripion	获取链码的总数  用于分页时计算页码数量
// @author	zhangbiao651
// @param	maps	interface{} 	链码信息的映射
// @param	count	int				链码信息的总数
// @param	err		error			查询过程中的错误信息
func GetChaincodeTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&ChaincodeInfo{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// @title GetChaincodeById
// @description	通过链码的 ID 获取链码信息
// @author	zhangbiao651
// @param	id		int		需要查询的链码的 ID
// @return	chaincode	*Chaincode	获取到的链码信息
func GetChaincode(name string) (*ChaincodeInfo, error) {
	var chaincode ChaincodeInfo
	err := db.Model(&ChaincodeInfo{}).Where("chaincode_id = ? ", name).First(&chaincode).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &chaincode, nil
}

// @title EditChaincode
// @description	更新 Chaincode 的信息
// @author 	zhangbiao651
// @param	id		int		链码的 id 用于数据库中查询链码
// @param	data	interface 更新后的信息
// @return	err		error	更新过程中的错误信息
func EditChaincode(id int, data interface{}) error {
	if err := db.Model(&ChaincodeInfo{}).Where("chaincode_id = ? ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// @title AddChaincode
// @description	向数据库中添加链码信息
// @author	zhangbiao651
// @param	data	map[string]interface{}	需要添加的链码信息
// @return	err		error		添加过程中发生的错误
func AddChaincode(data map[string]interface{}) error {
	chaincode := ChaincodeInfo{
		ChaincodeID:      data["chaincode_id"].(string),
		ChaincodePath:    data["chaincode_path"].(string),
		ChaincodeVersion: data["chaincode_version"].(string),
	}

	if err := db.Create(&chaincode).Error; err != nil {
		return err
	}
	return nil
}

// @title DeleteChaincode
// @description	删除链码
// @author	zhangbiao651
// @param	id	int		需要删除的链码的 id
// @return	err	error	删除过程中发生的错误
func DeleteChaincode(id int) error {
	if err := db.Where("chaincode_id = ?", id).Delete(&ChaincodeInfo{}).Error; err != nil {
		return err
	}
	return nil
}
