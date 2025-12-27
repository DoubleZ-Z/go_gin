package token

import (
	"go_gin/controllers"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	controllers.BaseController
}

// UploadFile 单文件上传
func (con *FileController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "上传失败",
		})
		return
	}
	// 上传文件到指定的目录
	dst := path.Join("./static/upload", file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "上传成功",
			"url":  "/static/upload/" + file.Filename,
		})
	}
}
