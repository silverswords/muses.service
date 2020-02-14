package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"muses.service/middleware"
	"muses.service/model"
)

func Login(ctx *gin.Context) {
	user := &model.User{}
	result := &model.Result{
		Code:    200,
		Message: "登录成功",
		Data:    "",
	}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	// 验证密码, 获取 user id
	if true {
		if token, err := middleware.JwtGenerateToken(user); err == nil {
			result.Message = "登录成功"
			result.Data = "Bearer " + token
			result.Code = http.StatusOK
			ctx.JSON(result.Code, gin.H{
				"result": result,
			})
		} else {
			result.Message = "登录失败"
			result.Code = http.StatusOK
			ctx.JSON(result.Code, gin.H{
				"result": result,
			})
		}
	} else {
		result.Message = "登录失败"
		result.Code = http.StatusOK
		ctx.JSON(result.Code, gin.H{
			"result": result,
		})
	}
}
