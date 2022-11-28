package models

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/jinzhu/gorm"
)

// fabric 网络的组织信息
type OrgInfo struct {
	gorm.Model
	// 组织的管理员用户名
	OrgAdminUser string `gorm:"column:org_admin;size:20;not null" `
	// 组织名称
	OrgName string `gorm:"column:org_name;size:255;not null;index:idx_ogn;unique" `
	// 组织的 MSPID
	OrgMspId string `gorm:"column:org_msp_id;size:20;not null" `
	// 组织的普通用户
	OrgUser string `gorm:"colunm:org_user;size:20;not null" `
	// 组织的节点个数
	OrgPeerNum int `gorm:"colunm:org_peer_num;not null" `

	// 组织锚节点文件的路径
	OrgAnchorFile string `gorm:"colunm:org_anchor_filepath;size:255;not null" `

	Peers []*PeerInfo `gorm:"ForeignKey:OrgName;AssociationForeignKey:OrgName"`
}

// @title ExistOrgInfoByOrgName
// @description 通过组织名判断是否存在组织
// @author zhangbiao651
// @param        orgname       string        组织名称
// @return       _			   bool 		 组织状态  true 表示该组织已经存在 false 表示该组织不存在
// @return       err 		   error		 执行过程中发生的错误
func ExistOrgInfoByOrgName(orgname string) (bool, error) {
	var orgInfo OrgInfo
	isExist := false
	log.Info("[models/org.go 开始查询：]")
	err := db.Select("id").Where("org_name = ? ", orgname).First(&orgInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if orgInfo.ID > 0 {
		isExist = true
	}
	log.Info("该组织当前状态为：", isExist)
	return isExist, nil
}

// @title GetOrgInfoTotal
// @decrtiption  返回数据库中的 OrgInfo 总数
// @author zhangbiao651
// @param		maps		interface{}		OrgInfo的信息
// @return		count		int				数据库中OrgInfo的总数
func GetOrgInfoTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&OrgInfo{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// @title GetOrgInfoByName
// @description 通过组织名获取组织信息
// @author zhangbiao651
// @param		orgname		string		需要获取组织信息的组织名
// @return 		orgInfo		*OrgInfo	获取到的组织信息
// @return		err			error		获取组织信息过程中出现的错误信息
func GetOrgInfoByName(orgname string) (*OrgInfo, error) {
	var orgInfo OrgInfo
	//err := db.Model(&OrgInfo{}).Where("org_name = ? ", orgname).Related(&orgInfo.Peers).First(&orgInfo).Error
	//err := db.Model(&OrgInfo{}).Where("org_name = ? ", orgname).First(&orgInfo).Error
	//err = db.Model(&orgInfo).Association("Peers").Error

	err := db.Model(&OrgInfo{}).Preload("Peers").Where("org_name = ?", orgname).First(&orgInfo).Error
	log.Info("[models/org.go/GetOrgInfoByName]获取到组织信息 %v", orgInfo)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &orgInfo, nil
}

// @title GetOrgInfos
// @description 分页查询组织信息
// @author zhangbiao651
// @param		pageNum			int			页面数
// @param		pageSize		int			页长
// @param		maps			interface{} OrgInfo信息
// @return		orgInfos		[]*Orginfo	该页所有组织的组织信息
func GetOrgInfos(pageNum int, pageSize int) ([]*OrgInfo, error) {
	var orgInfos []*OrgInfo
	err := db.Model(&OrgInfo{}).Preload("Peers").Limit(pageSize).Offset(pageNum).Find(&orgInfos).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return orgInfos, nil
}

// @title  AddOrgInfo
// @description  向数据库中添加组织信息
// @author zhangbiao651
// @param		orgInfo		OrgInfo		需要记录的组织信息
// @return 		err			error		记录组织信息时发生的错误
func AddOrgInfo(data map[string]interface{}) error {
	log.Info("models 开始获取组织信息")
	orgInfo := OrgInfo{
		OrgAdminUser:  data["org_admin"].(string),
		OrgAnchorFile: data["org_anchor_file"].(string),
		OrgName:       data["org_name"].(string),
		OrgUser:       data["org_user"].(string),
		OrgMspId:      data["org_msp_id"].(string),
		OrgPeerNum:    data["org_peer_num"].(int),
	}
	log.Info("组织信息获取完毕，开始向数据库写入组织信息")
	err := db.Create(&orgInfo).Error
	if err != nil {
		return err
	}
	err = db.Model(&orgInfo).Association("Peers").Error
	if err != nil {
		return err
	}
	return nil
}

// @title EditOrgInfo
// @description  将更新后的组织信息同步更新到数据库当中
// @author zhangbiao651
// @paraorgInfom		orgInfo		OrgInfo		需要同步到数据库的组织信息
// @return		err			error		更新组织信息时发生的错误
func EditOrgInfo(orgname string, data interface{}) error {
	log.Info("开始更新组织信息")
	if err := db.Model(&OrgInfo{}).Where("org_name = ? ", orgname).Update(data).Error; err != nil {
		return err
	}

	return nil
}

// @title DeleteOrgInfo
// @description  删除数据库中的组织信息
// @author zhangbiao651
// @param		orgname		string		需要删除的组织信息
// @return		err			error		删除数据过程中发生的错误
func DeleteOrgInfo(orgname string) error {
	if err := db.Where("org_name = ?", orgname).Delete(OrgInfo{}).Error; err != nil {
		return err
	}
	return nil
}
