package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/service/orderer_service"
)

type Orderer struct {
	OA  string `json:"orderer_admin"`
	OON string `json:"orderer_org_name"`
	OE  string `json:"orderer_endpoint"`
}

func AddOrderer(ctx *gin.Context) {
	var order Orderer
	ctx.ShouldBind(&order)
	service := orderer_service.OrdererService{
		OrdererAdminUser: order.OA,
		OrdererOrgName:   order.OON,
		OrdererEndpoint:  order.OE,
	}
	// 查询排序组织是否存在
	isExit, err := service.Exist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "查询发生错误，请重试！",
			"error":   err,
		})
		return
	}
	if isExit {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  2001,
			"message": "该排序组织名已经存在，请重命名新排序组织！",
		})
		return
	}

	// 存储排序组织信息
	err = service.Add()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "新增排序组织时发生错误，请重试",
			"error":   err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "添加排序组织成功",
	})
}

// @title GetOrdererByName
// @description 通过排序组织名获取排序组织信息
// @author	zhangbiao651
// @produce json
// @param	orgName
// @url	/Orderer/getOrdererbyname?org_name
func GetOrdererByName(ctx *gin.Context) {
	service := orderer_service.OrdererService{
		OrdererOrgName: ctx.Query("orderer_org_name"),
	}
	data, err := service.Get()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "查询是发生错误",
			"error":   err.Error(),
		})
		return
	}

	if data == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusNotFound,
			"message": "未找到该排序组织",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
}

// @title GetOrderers
// @description 查询所有排序组织的信息
// @url /Orderer/getall
func GetAllOrderer(ctx *gin.Context) {
	pageNum, _ := strconv.ParseInt(ctx.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 32)
	service := orderer_service.OrdererService{
		PageNum:  int(pageNum),
		PageSize: int(pageSize),
	}
	data, err := service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "获取排序失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
}
