package service

import "github.com/gin-gonic/gin"

func msg(err error, data gin.H) gin.H {
	if err == nil {
		return gin.H{
			"msg":  "success",
			"data": data,
		}
	} else {
		return gin.H{
			"msg":  err.Error(),
			"data": gin.H{},
		}
	}
}
