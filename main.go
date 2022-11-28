package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhangbiao651/fabric-manager/web/models"
	"github.com/zhangbiao651/fabric-manager/web/pkg/setting"
	"github.com/zhangbiao651/fabric-manager/web/routers"
)

func init() {
	setting.Setup()
	models.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routerInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routerInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server setUP errr :%v", err)
	}

}
