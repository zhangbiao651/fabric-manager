# 节点管理
## 添加节点
Url : http://127.0.0.1:8000/peer/add
### peer0.org1.zhangbiao651.top
{
	"org_name" : "Org1",
	"url" : "peer0.org1.zhangbiao651.top:7051"
}
### peer1.org1.zhangbiao651.top
{
	"org_name" : "Org1",
	"url" : "peer0.org1.zhangbiao651.top:9051"
}
### peer0.org2.zhangbiao651.top
{
	"org_name" : "Org2",
	"url" : "peer0.org2.zhangbiao651.top:7151"
}
### peer1.org2.zhangbiao651.top
{
	"org_name" : "Org2",
	"url" : "peer1.org2.zhangbiao651.top:9151"
}
## 删除节点
Url : http://127.0.0.1:8000/peer/delete?url=

# 组织管理
## 添加组织
Url : http://127.0.0.1:8000/org/add
### Org1

{
	"org_admin" : "Admin",
	"org_name" :  "Org1",
	"org_msp_id" : "Org1MSP",
	"org_user" : "User1",
	"org_anchor_file" : "/home/zb/fabric/project/v0.0.1/fabric-manager/fixtrues/data/",
	"org_peer_num": 2
}

### Org2
{
	"org_admin" : "Admin",
	"org_name" :  "Org2",
	"org_msp_id" : "Org2MSP",
	"org_user" : "User2",
	"org_anchor_file" : "/home/zb/fabric/project/v0.0.1/fabric-manager/fixtrues/data/",
	"org_peer_num": 2
}
## 通过组织名查询组织
Url : http://127.0.0.1:8000/org/get?org_name=
## 获取所有组织
Url : http://127.0.0.1:8000/org/get?page_name= && page_size=

# 排序节点管理
## 添加排序节点
Url ： http://127.0.0.1:8000/orderer/add
### orderer.zhangbiao651.top
{
	"orderer_admin" : "Admin",
	"orderer_org_name" : "OrdererOrg",
	"orderer_endpoint" : "orderer.zhangbiao651.top"
}
## 通过排序节点名称获取排序节点
Url ： http://127.0.0.1:8000/orderer/getbyname?orderer_org_name= 
## 获取所有排序节点
Url ： http://127.0.0.1:8000/orderer/getall?page_name= &&page_size=

# 通道信息管理
## 创建通道信息
Url : http://127.0.0.1:8000/channel/add
### mychannel
{
	"channel_id" :"mychannel",
	"channel_config_path":"/home/zb/fabric/project/v0.0.1/fabric-manager/fixtrues/data/mychannel.tx"
}
### testchannel
{
	"channel_id" :"mychannel",
	"channel_config_path":"/home/zb/fabric/project/v0.0.1/fabric-manager/fixtrues/data/mychannel.tx"
}

## 查询通道信息
### 通过通道ID获取通道信息
Url ： http://127.0.0.1:8000/channel/get?channel_id=
### 获取所有通道信息
Url： http://127.0.0.1:8000/channel/getall

## 向通道内添加组织
Url : http://127.0.0.1:8000/channel/orgjoin

### mychannel 
{
    "channel_id" : "mychannel",
    "orgnames":["Org1","Org2"]
}
### testchannel
{
    "channel_id" : "mychannel",
    "orgnames":["Org1","Org2"]
}




# 链码管理
## 添加链码
Url : http://127.0.0.1:8000/chaincode/add
### simplecc
{
	"chaincode_id" : "simplecc",
	"chaincode_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/chaincode/go/simplecc",
	"chaincode_version":"1.0.0"
}

## 通过链码名称获取链码
http://127.0.0.1:8000/chaincode/get?chaincode_id=

## 获取所有链码
Url: http://127.0.0.1:8000/chaincode/getall?pagenum=0&pagesize=20

# 事件信息
## 创建通道
URL ： http://127.0.0.1:8000/env/channelcreate
### mychannel 
{
    "channel_id" : "mychannel",
    "orderer_name" : "OrdererOrg",
    "peers":["peer0.org1.zhangbiao651.top:7051","peer1.org1.zhangbiao651.top:9051","peer0.org2.zhangbiao651.top:7151","peer1.org2.zhangbiao651.top:9151"],
    "chaincode_id" : "simplecc",
    "config_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
}
### testchannel
{
    "channel_id" : "testchannel",
    "orderer_name" : "OrdererOrg",
    "peers":["peer0.org1.zhangbiao651.top:7051","peer1.org1.zhangbiao651.top:9051","peer0.org2.zhangbiao651.top:7151","peer1.org2.zhangbiao651.top:9151"],
    "chaincode_id" : "simplecc",
    "config_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
}

## 节点加入通道
URL ： http://127.0.0.1:8000/env/joinchannel
### mychannel
{
    "channel_id" : "mychannel",
    "orderer_name" : "OrdererOrg",
    "peers":["peer0.org1.zhangbiao651.top:7051","peer1.org1.zhangbiao651.top:9051","peer0.org2.zhangbiao651.top:7151","peer1.org2.zhangbiao651.top:9151"],
    "chaincode_id" : "simplecc",
    "config_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
}
### testchannel
{
    "channel_id" : "testchannel",
    "orderer_name" : "OrdererOrg",
    "peers":["peer0.org1.zhangbiao651.top:7051","peer0.org2.zhangbiao651.top:7151"],
    "chaincode_id" : "simplecc",
    "config_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
}

## 安装链码
URL ： http://127.0.0.1:8000/env/createcc
### mychannel
{
    "channel_id" : "mychannel",
    "orderer_name" : "OrdererOrg",
    "peers":["peer0.org1.zhangbiao651.top:7051","peer1.org1.zhangbiao651.top:9051","peer0.org2.zhangbiao651.top:7151","peer1.org2.zhangbiao651.top:9151"],
    "chaincode_id" : "simplecc",
    "config_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
}
### testchannel
{
    "channel_id" : "testchannel",
    "orderer_name" : "OrdererOrg",
    "peers":["peer0.org1.zhangbiao651.top:7051","peer0.org2.zhangbiao651.top:7151"],
    "chaincode_id" : "simplecc",
    "config_path" :"/home/zb/fabric/project/v0.0.1/fabric-manager/config.yaml"
}


