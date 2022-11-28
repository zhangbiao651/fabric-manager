package api

import (
	"net/http"
	"strconv"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/service/channel_service"
)

type OACData struct {
	Channel  string   `json:"channel_id"`
	OrgNames []string `json:"orgnames"`
}

type Channel struct {
	ChannelID string `json:"channel_id"`
	ChannelConfigPath string `json:"channel_config_path"`
}

func AddChannel(ctx *gin.Context) {
	var channel Channel
	ctx.ShouldBind(&channel)
	service := channel_service.ChannelService{
		ChannelID:         channel.ChannelID,
		ChannelConfigPath: channel.ChannelConfigPath,
	}
	log.Info("[api/channel.go]channel-id is", service.ChannelID)
	isExist, err := service.Exist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "查看通道信息失败",
			"error":   err.Error(),
		})
		return
	}

	if isExist {
		ctx.JSON(http.StatusOK, gin.H{
			"statsu":  2001,
			"message": "该通道已经存在",
		})
		return
	}

	err = service.Add()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "添加失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "通道添加成功",
	})
}

func GetChannels(ctx *gin.Context) {
	pageNum, _ := strconv.ParseInt(ctx.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 32)
	service := channel_service.ChannelService{
		PageNum:  int(pageNum),
		PageSize: int(pageSize),
	}
	data, err := service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "获取通道列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "获取成功",
		"data":    data,
	})
}

func GetChannel(ctx *gin.Context) {
	service := channel_service.ChannelService{
		ChannelID: ctx.Query("channel_id"),
	}
	data, err := service.Get()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "获取通道列表失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "获取成功",
		"data":    data,
	})

}

func OrgJoinChannel(ctx *gin.Context) {
	var oacd OACData
	ctx.ShouldBind(&oacd)
	log.Info("[api OrgJoinChannel] 获取到参数 %v", oacd)
	service := &channel_service.ChannelService{
		ChannelID: oacd.Channel,
		OrgNames:  oacd.OrgNames,
	}
	err := service.AddOrg()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"message ": "组织加入通道失败",
			"error":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "已将所有组织添加到通道",
	})

}
func OrgLeaveChannel(ctx *gin.Context) {
	var oacd OACData
	ctx.ShouldBind(&oacd)
	log.Info("[api OrgJoinChannel] 获取到参数 %v", oacd)
	service := &channel_service.ChannelService{
		ChannelID: oacd.Channel,
		OrgNames:  oacd.OrgNames,
	}
	err := service.MoveOrg()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"message ": "组织移出通道失败",
			"error":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "已将组织移出通道",
	})

}
