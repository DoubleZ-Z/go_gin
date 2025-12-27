package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (con *BaseController) Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
	})
}

func (con *BaseController) Fail(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": "失败",
	})
}
