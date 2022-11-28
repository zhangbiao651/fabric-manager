package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/service/chaincode_service"
)

type chaincode struct {
	ChaincodeID      string `json:"chaincode_id"`
	ChaincodePath    string `json:"chaincode_path"`
	ChaincodeVersion string `json:"chaincode_version"`
}

func AddChaincode(ctx *gin.Context) {
	var code chaincode
	ctx.ShouldBind(&code)
	service := chaincode_service.ChaincodeService{
		ChaincodeID:      code.ChaincodeID,
		ChaincodePath:    code.ChaincodePath,
		ChaincodeVersion: code.ChaincodeVersion,
	}
	err := service.Add()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "添加链码失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "添加链码成功",
	})
}

func GetChaincode(ctx *gin.Context) {
	chaincodeID := ctx.Query("chaincode_id")
	service := chaincode_service.ChaincodeService{
		ChaincodeID: chaincodeID,
	}
	data, err := service.Get()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "链码信息获取失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	return
}
func GetAllChaincode(ctx *gin.Context) {
	pageNum, _ := strconv.ParseInt(ctx.Query("page_num"), 10, 32)
	pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 32)
	service := chaincode_service.ChaincodeService{
		PageNum:  int(pageNum),
		PageSize: int(pageSize),
	}
	data, err := service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "链码信息获取失败",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	return
}
