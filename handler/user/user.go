package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"muses.service/middleware"
	usermodel "muses.service/model/user"
)

type Controller struct {
	db *sql.DB
}

// New create an external service interface
func NewDB(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	name := "Admin"
	password := "111111"
	err := usermodel.CreateTable(c.db, &name, &password)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/register", c.create)
	r.POST("/login", c.login)
	r.POST("/sendMsg", c.sendMsg)
}

func (c *Controller) create(ctx *gin.Context) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	//Default password
	if admin.Password == "" {
		admin.Password = "111111"
	}

	err = usermodel.Create(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) login(ctx *gin.Context) {
	user := &usermodel.User{}
	result := &usermodel.Result{
		Code:    200,
		Message: "登录成功",
		Data:    "",
	}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	_, err := usermodel.Login(c.db, &user.Name, &user.Password)
	if err == nil {
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

func (c *Controller) sendMsg(ctx *gin.Context) {

}
