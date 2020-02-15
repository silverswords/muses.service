package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"muses.service/model"
)

type Room struct {
	Id string `json:"id"`
}

func Joinroom(ctx *gin.Context) {
	result := &model.Result{}
	room := &Room{}
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	result.Message = "操作成功"
	result.Code = http.StatusOK
	ctx.JSON(result.Code, gin.H{
		"result": result,
	})
}
