package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/service/env_service"
)

type EnvInfo struct {
	ChannelID   string   `json:"channel_id"`
	OrdererName string   `json:"orderer_name"`
	Peers       []string `json:"peers"`
	ChaincodeID string   `json:"chaincode_id"`
	ConfigPath  string   `json:"config_path"`
}

func CreateChannel(ctx *gin.Context) {
	var envInfo EnvInfo
	ctx.ShouldBindJSON(&envInfo)
	envService := &env_service.EnvService{
		ChannelID:   envInfo.ChannelID,
		OrdererName: envInfo.OrdererName,
		ChaincodeID: envInfo.ChaincodeID,
		PeerUrl:     envInfo.Peers,
		ConfigPath:  envInfo.ConfigPath,
	}
	err := envService.Init()

	log.Printf("[api/env.go/CreateAndJoinChannel]info.OrdererClientContext =[%v] ", envService.SdkInfo.OrdererClientContext)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "创建通道失败",
			"error":   err.Error(),
		})
		return
	}
	if err := envService.CreateChannel(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "创建通道失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "已完成通道创建",
	})
}

func JoinChannel(ctx *gin.Context) {
	var envInfo EnvInfo
	ctx.ShouldBindJSON(&envInfo)
	envService := &env_service.EnvService{
		ChannelID:   envInfo.ChannelID,
		OrdererName: envInfo.OrdererName,
		ChaincodeID: envInfo.ChaincodeID,
		PeerUrl:     envInfo.Peers,
		ConfigPath:  envInfo.ConfigPath,
	}
	log.Print("[INFO][api/env.go/JoinChannel] ConfigPath =  ", envInfo.ConfigPath)
	if err := envService.Init(); err != nil {
		envService.Sdk.Close()
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "初始化环境失败",
			"error":   err.Error(),
		})
		return
	}

	if err := envService.JoinChannel(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "节点加入失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "已将节点加入到通道中",
	})
	return
}

func CreateCCLifecycle(ctx *gin.Context) {
	var envInfo EnvInfo

	ctx.ShouldBindJSON(&envInfo)
	envService := &env_service.EnvService{
		ChannelID:   envInfo.ChannelID,
		OrdererName: envInfo.OrdererName,
		ChaincodeID: envInfo.ChaincodeID,
		PeerUrl:     envInfo.Peers,
		ConfigPath:  envInfo.ConfigPath,
	}

	if err := envService.Init(); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "环境初始化失败",
			"error":   err.Error(),
		})
		return
	}
	if err := envService.CreateCCLifecycle(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "链码安装失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "链码安装成功",
	})
}
