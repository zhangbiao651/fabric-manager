package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	apiPeer := r.Group("/peer")
	{
		apiPeer.POST("/add", api.AddPeer)
		apiPeer.GET("/getall", api.GetPeers)
		apiPeer.GET("/delete",api.DeletePeer)

	}
	apiOrg := r.Group("/org")
	{
		apiOrg.POST("/add", api.AddOrgInfo)
		apiOrg.GET("/get", api.GetOrgInfoByName)
		apiOrg.GET("/getall", api.GetAllOrg)

	}
	apiOrderer := r.Group("/orderer")
	{
		apiOrderer.POST("/add", api.AddOrderer)
		apiOrderer.GET("/getbyname", api.GetOrdererByName)
		apiOrderer.GET("/getall", api.GetAllOrderer)
	}
	apiEnv := r.Group("/env")
	{
		apiEnv.POST("/channelcreate", api.CreateChannel)
		apiEnv.POST("/joinchannel", api.JoinChannel)
		apiEnv.POST("/createcc", api.CreateCCLifecycle)
	}
	apiChannel := r.Group("channel")
	{
		apiChannel.POST("/add", api.AddChannel)
		apiChannel.GET("/get", api.GetChannel)
		apiChannel.GET("/getall", api.GetChannels)
		apiChannel.POST("/orgjoin", api.OrgJoinChannel)
		apiChannel.POST("/orgleave",api.OrgLeaveChannel)
	}

	apiChaincode := r.Group("/chaincode")
	{
		apiChaincode.POST("/add", api.AddChaincode)
		apiChaincode.GET("/get", api.GetChaincode)
		apiChaincode.GET("/getall", api.GetAllChaincode)
	}

	return r
}
