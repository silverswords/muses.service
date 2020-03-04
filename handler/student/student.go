package student

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gorp.v1"
	"muses.service/middleware"
)

// Controller -
type Controller struct {
	dbmap *gorp.DbMap
}

// Student -
type Student struct {
	UserID   string
	Created  int64
	Name     string
	Password string
	Role     string
	IsBusy   bool
}

// NewDB -
func NewDB(dbmap *gorp.DbMap) *Controller {
	return &Controller{
		dbmap: dbmap,
	}
}

// RegisterRouter -
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	c.dbmap.AddTableWithName(Student{}, "student").SetKeys(false, "UserID")

	fmt.Print("create person")

	r.POST("/create", c.create)
	r.POST("/login", c.login)
	r.GET("/tourist", c.createTourist)
	r.POST("/changePassword", c.changePassword)
	r.POST("/changename", c.changeName)
	r.POST("/sendMsg", c.sendMsg)
}

func (c *Controller) create(ctx *gin.Context) {
	var (
		user struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
			Role     string `json:"role"`
			IsBusy   bool   `json:"isbusy"`
		}
	)

	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if user.Password == "" && user.Name == "" {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	student := Student{
		UserID:   uuid.NewV4().String(),
		Created:  time.Now().UnixNano(),
		Name:     user.Name,
		Password: user.Password,
		Role:     user.Role,
		IsBusy:   user.IsBusy,
	}

	err = c.dbmap.Insert(&student)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (c *Controller) createTourist(ctx *gin.Context) {
	student := Student{
		UserID:   uuid.NewV4().String(),
		Created:  time.Now().UnixNano(),
		Name:     "XXX",
		Password: "",
		Role:     "student",
		IsBusy:   false,
	}

	err := c.dbmap.Insert(&student)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   student,
	})
}

// Login -
func (c *Controller) login(ctx *gin.Context) {
	var (
		user struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
			Role     string `json:"role"`
		}
	)

	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	person := Student{
		Name:     user.Name,
		Password: user.Password,
		Role:     user.Role,
	}

	err = c.dbmap.SelectOne(&person, "select * from person where name=? and password=? and role = ? limit 1", person.Name, person.Password, person.Role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	token, err := middleware.JwtGenerateToken(person.UserID.String())
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   token,
	})
}

func (c *Controller) changeName(ctx *gin.Context) {
	var (
		user struct {
			ID   string `json:"id"`
			Name string `json:"name"      binding:"required,alphanum,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	person := Student{
		UserID: user.ID,
		Name: user.Name,
	}

	err := c.dbmap.Update(&person)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK
	})
}

func (c *Controller) changePassword(ctx *gin.Context) {
	var (
		user struct {
			ID   string `json:"id"`
			Password string `json:"password"      binding:"required,alphanum,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	person := Student{
		UserID: user.ID,
		Password: user.Password,
	}

	err := c.dbmap.Update(&person)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK
	})
}

func (c *Controller) sendMsg(ctx *gin.Context) {

}
