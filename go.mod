module github.com/zhangbiao651/fabric-manager

go 1.14

replace (
	github.com/zhangbiao651/fabric-manager/fabtools => ./fabtools
	github.com/zhangbiao651/fabric-manager/fabtools/fabric => ./fabtools/fabric
	github.com/zhangbiao651/fabric-manager/fabtools/fabric_ca => ./fabtools/fabric_ca
	github.com/zhangbiao651/fabric-manager/fabtools/pkg/cryptogen/ca => ./fabtools/pkg/cryptogen/ca
	github.com/zhangbiao651/fabric-manager/fabtools/pkg/cryptogen/csp => ./fabtools/pkg/cryptogen/csp
	github.com/zhangbiao651/fabric-manager/sdkInit => ./sdkInit
	github.com/zhangbiao651/fabric-manager/web/api => ./web/api
	github.com/zhangbiao651/fabric-manager/web/models => ./web/models
	github.com/zhangbiao651/fabric-manager/web/pkg/file => ./web/pkg/file
	github.com/zhangbiao651/fabric-manager/web/pkg/logging => ./web/pkg/logging
	github.com/zhangbiao651/fabric-manager/web/pkg/setting => ./web/pkg/setting
	github.com/zhangbiao651/fabric-manager/web/routers => ./web/routers
	github.com/zhangbiao651/fabric-manager/web/service => ./web/service
	github.com/zhangbiao651/fabric-manager/web/service/chaincode_service => ./web/service/chaincode_service
	github.com/zhangbiao651/fabric-manager/web/service/channel_service => ./web/service/channel_service
	github.com/zhangbiao651/fabric-manager/web/service/env_service => ./web/service/env_service
	github.com/zhangbiao651/fabric-manager/web/service/orderer_service => ./web/service/orderer_service
	github.com/zhangbiao651/fabric-manager/web/service/org_service => ./web/service/org_service
	github.com/zhangbiao651/fabric-manager/web/service/peer_service => ./web/service/peer_service

)

require (
	github.com/cloudflare/cfssl v1.6.3
	github.com/gin-gonic/gin v1.8.1
	github.com/golang/protobuf v1.5.2
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-amcl v0.0.0-20221107192335-5c75bc7be9c0 // indirect
	github.com/hyperledger/fabric-config v0.1.0
	github.com/hyperledger/fabric-protos-go v0.0.0-20210318103044-13fdee960194
	github.com/hyperledger/fabric-sdk-go v1.0.1-0.20210603143513-14047c6d88f0
	github.com/jinzhu/gorm v1.9.16
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	gopkg.in/ini.v1 v1.67.0
	gorm.io/gorm v1.24.2

)
