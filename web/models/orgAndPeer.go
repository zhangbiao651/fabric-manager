package models

type EnvInfo struct {
	ChannelID   string   `json:"channelID"`
	OrdererName string   `json:"orderer_name"`
	OrgNames    []string `json:"org_names"`
	PeerUrl     []string `json:"peer_url"`
	ChaincodeID string   `json:"chaincode_id"`
	ConfigPath  string   `json:"config_path"`
}
