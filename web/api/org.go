package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/service/org_service"
)

// @title SerOrgInfo
// @description 添加组织信息到数据库中
// @author zhangbiao651
// @produce json
// @url : /orginfo/addorgInfo
type OgnInfo struct {
	OrgAdminUser  string `json:"org_admin"`
	OrgName       string `json:"org_name"`
	OrgMspId      string `json:"org_msp_id"`
	OrgUser       string `json:"org_user"`
	OrgAnchorFile string `json:"org_anchor_file"`
	OrgPeerNum    int    `json:"org_peer_num"`
}

func AddOrgInfo(ctx *gin.Context) {
	log.Printf("[INFO][api/org.go/AddOrgInfo] 开始添加组织")
	var ogn OgnInfo
	ctx.ShouldBind(&ogn)
	service := org_service.OrgService{
		OrgAdminUser:  ogn.OrgAdminUser,
		OrgName:       ogn.OrgName,
		OrgMspId:      ogn.OrgMspId,
		OrgUser:       ogn.OrgUser,
		OrgAnchorFile: ogn.OrgAnchorFile,
		OrgPeerNum:    ogn.OrgPeerNum,
	}
	// 查询组织是否存在
	isExit, err := service.Exist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "查询发生错误，请重试！",
			"error":   err.Error(),
		})
		return
	}
	if isExit {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  2001,
			"message": "该组织已存在",
		})
		return

	}

	// 存储组织信息
	err = service.Add()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "新增组织时发生错误，请重试",
			"error":   err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "添加组织成功",
	})
}

// @title GetOrgInfoByName
// @description 通过组织名获取组织信息
// @author	zhangbiao651
// @produce json
// @param	orgName
// @url	/orgInfo/getorginfobyname?org_name
func GetOrgInfoByName(ctx *gin.Context) {
	service := org_service.OrgService{
		OrgName: ctx.Query("org_name"),
	}
	log.Printf("[info][api/getorgInfoByName] service  获取到组织名%s", ctx.Query("org_name"))
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
			"message": "未找到该组织",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
}

// @title GetOrgInfos
// @description 查询所有组织的信息
// @url /orgInfo/getall
func GetAllOrg(ctx *gin.Context) {
	pageNum, _ := strconv.ParseInt(ctx.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 32)
	service := org_service.OrgService{
		PageNum:  int(pageNum),
		PageSize: int(pageSize),
	}
	data, err := service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "获取组织失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
}

func EditOrgInfo(ctx *gin.Context) {
	log.Printf("[info] 开始更新组织")
	var ogn OgnInfo
	ctx.ShouldBind(&ogn)
	service := org_service.OrgService{
		OrgAdminUser:  ogn.OrgAdminUser,
		OrgName:       ogn.OrgName,
		OrgMspId:      ogn.OrgMspId,
		OrgUser:       ogn.OrgUser,
		OrgAnchorFile: ogn.OrgAnchorFile,
		OrgPeerNum:    ogn.OrgPeerNum,
	}

	// 查询组织是否存在
	isExit, err := service.Exist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "查询发生错误，请重试！",
			"error":   err.Error(),
		})
		return
	}
	if isExit {
		err = service.Edit()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "更新组织信息时发生错误",
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "已更新组织信息",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusInternalServerError,
		"message": "该组织不存在",
	})
	return

}
