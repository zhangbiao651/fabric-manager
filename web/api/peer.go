package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/models"
	"github.com/zhangbiao651/fabric-manager/web/service/peer_service"
)

type Data struct {
	Peers []*models.PeerInfo `json:"peers"`
}

type Peer struct {
	OrgName string `json:"org_name"`
	Url     string `json:"url"`
}

// @description  添加节点
// @url  post  /peer/addpeer
func AddPeer(ctx *gin.Context) {
	var peer Peer
	ctx.ShouldBind(&peer)
	service := peer_service.PeerService{
		OrgName: peer.OrgName,
		Url:     peer.Url,
	}
	log.Printf("[info][api/peer.go/Addpeer],开始[%s]节点的查重", service.Url)
	isExist, err := service.Exist()

	log.Printf("[info][api/peer.go/Addpeer],[%s]节点的查重结果为[%v]", service.Url, isExist)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "peer 的查重失败，请重试",
			"error":   err,
		})
		return
	}
	if isExist {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  2001,
			"message": "该节点已经存在",
		})
		return
	}
	err = service.Add()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "peer 节点添加失败",
			"error":   err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "节点添加成功",
	})

}

// getAll
func GetPeers(ctx *gin.Context) {
	pageNum, _ := strconv.ParseInt(ctx.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 32)

	service := peer_service.PeerService{
		OrgName:  ctx.Query("org_name"),
		PageNum:  int(pageNum),
		PageSize: int(pageSize),
	}
	log.Printf("[info] orgName is %s", service.OrgName)
	peers, err := service.GetAll()
	//	data := Data{
	//		Peers: peers,
	//	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "查找节点信息失败！",
			"error":   err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "查找成功",
		"data":    peers,
	})
}

func DeletePeer(ctx *gin.Context) {
	service := peer_service.PeerService{
		Url: ctx.Query("url"),
	}
	isExist, err := service.Exist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "peer 的查重失败，请重试",
			"error":   err,
		})
		return

	}
	if !isExist {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "要删除的节点不存在,请重试",
			"error":   err.Error(),
		})
		return
	}
	if err := service.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "删除节点失败,请重试",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "删除节点成功",
	})
}
