package router

import (
	"github.com/ciiim/cloudborad/cmd/backpack/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery(), cors.Default())
	initRouter(g)
	return g
}

func initRouter(g *gin.Engine) {

	apiv1Group := g.Group("/api/v1")
	{
		fileGroup := apiv1Group.Group("/file")
		{
			fileGroup.POST("/info/:space/*path", service.GetFileInfo)

			// 检查文件块是否存在
			fileGroup.POST("/check/chunk", service.CheckFileChunk)

			// 获取文件内容
			fileGroup.GET("/:space/*path", service.GetFileContent)

			// 上传文件
			fileGroup.POST("/:space/upload", service.UploadFile)

			// fileGroup.POST("/:space/upload/chunk", service.UploadFileChunk)

			//删除文件
			fileGroup.DELETE("/:space/*path", service.DeleteFile)
		}

		spaceNoAuthGroup := apiv1Group.Group("/space/noauth")
		{
			// 创建space
			spaceNoAuthGroup.POST("/", service.CreateSpace)

			spaceNoAuthGroup.POST("/verify", service.VerifySpace)

			spaceNoAuthGroup.POST("/list", service.ListSpaces)
		}
		spaceGroup := apiv1Group.Group("/space", service.SpaceTokenAuth)
		{

			spaceGroup.POST("/:space/stat", service.GetSpaceStat)

			spaceGroup.POST("/:space/modify", service.SetSpaceStat)

			spaceGroup.POST("/:space/create_dir", service.CreateDir)

			spaceGroup.POST("/:space/delete_dir", service.DeleteDir)

			spaceGroup.POST("/:space/rename_dir", service.RenameDir)

			spaceGroup.POST("/:space/delete", service.DeleteSpace)

			spaceGroup.POST("/:space/get", service.GetDirInSpace)

		}

		serverGroup := apiv1Group.Group("/server")
		{
			serverGroup.POST("/join", service.JoinServer)

			serverGroup.POST("/leave", service.LeaveServer)

			serverGroup.POST("/list", service.ListServer)

			serverGroup.POST("/info", service.ServerInfo)

			serverGroup.GET("/feature", service.GetFeatureList)

			serverGroup.POST("/feature", service.FeatureControl)

		}
	}

}
