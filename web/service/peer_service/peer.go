package peer_service

import (
	"log"

	"github.com/zhangbiao651/fabric-manager/web/models"
)

type PeerService struct {
	OrgName string
	Url     string

	PageSize int
	PageNum  int
}

func (service *PeerService) Exist() (bool, error) {
	isExist, err := models.ExistPeerByUrl(service.Url)
	if err != nil {
		return false, err
	}
	if isExist {
		log.Printf("[peer_service/peer.go/Exist],[%s]已经存在", service.Url)
	}
	return isExist, nil
}

func (service *PeerService) GetAll() ([]*models.PeerInfo, error) {
	peers, err := models.GetPeers(service.PageNum, service.PageSize, service.OrgName)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (service *PeerService) Get() (*models.PeerInfo, error) {
	peer, err := models.GetPeerByUrl(service.Url)
	if err != nil {
		return nil, err
	}
	return peer, nil
}

func (service *PeerService) Add() error {
	peer := map[string]interface{}{
		"org_name": service.OrgName,
		"url":      service.Url,
	}

	err := models.AddPeer(peer)
	if err != nil {
		return err
	}

	return nil
}

func (service *PeerService) Delete() error {
	err := models.DeletePeer(service.Url)
	return err
}
