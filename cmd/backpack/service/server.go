package service

import (
	"net/http"

	"github.com/ciiim/cloudborad/cmd/backpack/model"
	"github.com/gin-gonic/gin"
)

func JoinServer(ctx *gin.Context) {
	addr := ctx.PostForm("addr")
	if addr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": ErrParam})
		return
	}
	if err := model.Ring.NodeService.Join(addr); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, msg(nil, nil))
}

func LeaveServer(ctx *gin.Context) {
	model.Ring.NodeService.Shutdown()
	ctx.JSON(http.StatusOK, msg(nil, nil))
}

func ListServer(ctx *gin.Context) {
	nodes := model.Ring.NodeService.NodeServiceRO().GetAllReal()

	ctx.JSON(http.StatusOK, msg(nil, gin.H{"nodes": nodes}))
}

func ServerInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, msg(nil, gin.H{"info": model.Ring.NodeService.NodeServiceRO().Self()}))
}

func GetFeatureList(ctx *gin.Context) {
	//TODO:完成模块系统后补充
}

func FeatureControl(ctx *gin.Context) {
	feature := ctx.PostForm("feature")
	enable := ctx.PostForm("enable")

	if feature == "" || enable == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": ErrParam})
		return
	}
}
