package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type PeerInfo struct {
	gorm.Model
	OrgName string `json:"org_name" gorm:"not null"`

	Url string `json:"url" gorm:"not null;unique"`
}

func ExistPeerByUrl(url string) (bool, error) {
	var peer PeerInfo
	err := db.Select("id").Where("url = ?", url).First(&peer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	log.Println("[info][models/peer.go]查询到的peer的ID 为", peer.ID)
	if peer.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetPeerTotal(orgname string) (int, error) {
	var count int
	if err := db.Model(&PeerInfo{}).Where("org_name = ", orgname).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetPeers(pageNum int, pageSize int, orgname string) ([]*PeerInfo, error) {
	log.Printf("[info] GetPeers by orgname=%s", orgname)
	var peers []*PeerInfo
	//err := db.Where("org_name = ?", orgname).Offset(pageNum).Limit(pageSize).Find(&peers).Error
	err := db.Where("org_name = ?", orgname).Offset(pageNum).Limit(pageSize).Find(&peers).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	log.Printf("[info] peers is %v", peers)

	return peers, nil
}

func GetPeerByUrl(url string) (*PeerInfo, error) {
	var peer PeerInfo
	err := db.Where("url = ?", url).First(&peer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &peer, nil
}

func AddPeer(data map[string]interface{}) error {
	peer := PeerInfo{
		OrgName: data["org_name"].(string),
		Url:     data["url"].(string),
	}
	log.Println("[models/peer.go] the peer Info is ", peer)
	// 添加节点与组织的关联

	if err := db.Create(&peer).Error; err != nil {
		return err
	}
	//peers, err := GetPeerByUrl(peer.Url)
	//if err != nil {
	//	return err
	//}
	//orgInfo, err := GetOrgInfoByName(peer.OrgName)
	//err = db.Model(&orgInfo).Association("Peers").Append(peers).Error
	//if err != nil {
	//	log.Printf("[models/peer.go/AddPeer] 创建节点与组织关联时发生错误%s", err.Error())
	//}
	return nil
}

func DeletePeer(url string) error {
	peer, err := GetPeerByUrl(url)
	if err != nil {
		return fmt.Errorf("[models/peer.go/DeletePeer]查询节点失败")
	}
	err = db.Delete(peer).Error
	if err != nil {
		return err
	}
	return nil
}
