package apis

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

}

func apiGrp(r *gin.Engine) {
	apiGrp := r.Group("/api/v1")

	userGrp := apiGrp.Group("/user")
	{
		userGrp.POST("/default")
		userGrp.POST("/login")
	}
}
