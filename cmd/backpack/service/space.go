package service

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ciiim/cloudborad/cmd/backpack/model"
	"github.com/gin-gonic/gin"
)

type BaseAndDir struct {
	Token string `json:"token"`
	Base  string `json:"base"`
	Dir   string `json:"dir"`
}

/*
	 JSON Body
		{
			"token": "<token>",
			"base": "<base>",
			"dir": "<dir>"
		}
*/
func GetDirInSpace(ctx *gin.Context) {
	space := ctx.Param("space")
	var baseAndDir BaseAndDir
	if err := ctx.BindJSON(&baseAndDir); err != nil {
		println(err.Error())
		ctx.JSON(http.StatusOK, msg(ErrParam, nil))
		return
	}

	userDir, err := model.Ring.DirInSpace(space, baseAndDir.Base, baseAndDir.Dir)
	if err != nil {
		ctx.JSON(http.StatusOK, msg(err, nil))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": userDir,
	})
}

/*
	 JSON Body
		{
			"space": "<space name>",
			"password": "<password>"
		}
*/
type createSpaceJsonBody struct {
	Space    string `json:"space"`
	Password string `json:"password"`
}

func CreateSpace(ctx *gin.Context) {
	var body createSpaceJsonBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": ErrParam})
		return
	}
	if body.Password == "" {
		if err := CreateSpaceNoPassword(body.Space); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
			return
		}
	} else {
		if err := CreateSpaceWithPassword(body.Space, body.Password); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
			return
		}
	}
	token := model.EncryptPassword(body.Password)
	println(token)
	ctx.JSON(http.StatusOK, msg(nil, gin.H{
		"space": body.Space,
		"token": token,
	}))
}

func CreateSpaceNoPassword(space string) error {
	return model.Ring.NewSpace(space, 0)
}

func CreateSpaceWithPassword(space, password string) error {
	if model.Ring.NewSpace(space, 0) != nil {
		return ErrSpaceExist
	}
	return model.Ring.SetSpacePassword(space, password)
}

func DeleteSpace(ctx *gin.Context) {
	space := ctx.PostForm("space")
	password := ctx.PostForm("password")

	if !model.Ring.CheckSpacePassword(space, password) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "password error"})
		return
	}

	if err := model.Ring.FrontSystem.DeleteSpace(space); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})

}

func GetSpaceStat(ctx *gin.Context) {
	space := ctx.Param("space")
	key := ctx.PostForm("key")
	stat, err := model.Ring.GetSpaceStat(space, key)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, msg(nil, gin.H{"stat": stat}))
}

func SetSpaceStat(ctx *gin.Context) {
	space := ctx.Param("space")
	key := ctx.PostForm("key")
	value := ctx.PostForm("value")
	if err := model.Ring.SetSpaceStat(space, key, value); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func CreateDir(ctx *gin.Context) {
	space := ctx.PostForm("space")
	base := ctx.PostForm("base")
	dir := ctx.PostForm("dir")
	err := model.Ring.MakeDir(space, base, dir)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func DeleteDir(ctx *gin.Context) {
	space := ctx.PostForm("space")
	base := ctx.PostForm("base")
	dir := ctx.PostForm("dir")
	err := model.Ring.DeleteDir(space, base, dir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
}

func RenameDir(ctx *gin.Context) {
	space := ctx.PostForm("space")
	base := ctx.PostForm("base")
	dir := ctx.PostForm("dir")
	newDirName := ctx.PostForm("new_dir_name")
	err := model.Ring.RenameDir(space, base, dir, newDirName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
}

func ListSpaces(ctx *gin.Context) {
	lists := model.Ring.GetAllSpaces()
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "data": lists})
}

type verifyBody struct {
	Space    string `json:"space"`
	Password string `json:"password"`
}

func VerifySpace(ctx *gin.Context) {
	var body verifyBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, msg(ErrParam, nil))
		return
	}

	if !model.Ring.CheckSpacePassword(body.Space, body.Password) {
		ctx.JSON(http.StatusOK, msg(ErrPassword, nil))
		return
	}
	ctx.JSON(http.StatusOK, msg(nil, gin.H{
		"space": body.Space,
		"token": model.EncryptPassword(body.Password),
	}))
}

type authBody struct {
	Token string `json:"token"`
}

func SpaceTokenAuth(ctx *gin.Context) {
	var body authBody
	readBodyAndSetBodyRepeatRead(ctx, func() {
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusOK, msg(ErrParam, nil))
			return
		}
	})
	space := ctx.Param("space")

	if !model.Ring.CheckSpaceToken(space, body.Token) {
		ctx.JSON(http.StatusOK, msg(ErrPassword, nil))
		return
	}
	ctx.Next()
}

func readBodyAndSetBodyRepeatRead(c *gin.Context, cb func()) {
	if s, ok := c.Request.Body.(io.Seeker); ok {
		//执行读取Body的操作
		cb()
		//再次设置可读状态
		_, err := s.Seek(0, 0)
		if err == nil {
			return
		}
	}

	bs, _ := io.ReadAll(c.Request.Body)
	_ = c.Request.Body.Close() // NOTE 原始的 Body 无需手动关闭,会在 response.reqBody中自动关闭的.
	//设置可读状态
	r := bytes.NewReader(bs)
	c.Request.Body = io.NopCloser(r)
	//执行读取Body的操作
	cb()
	//再次设置可读状态
	_, _ = r.Seek(0, 0)
}
