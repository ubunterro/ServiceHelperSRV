// Сервер для работы с приложением ServiceHelper
package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(headersByRequestURI())

	//r.Static("/public", "./public")

	r.LoadHTMLGlob("templates/*")

	admin := r.Group("/public", basicAuth())
	{
		admin.Static("/assets", "./public/assets")
		admin.GET("/index.html", Homepage)
	}

	client := r.Group("/api", basicAuth())
	{

		client.GET("/auth", Auth)

		client.GET("/api", Read)
		client.POST("/create", Create)
		client.POST("/edit", Edit)
		client.POST("/delete/:id", Delete)

		client.GET("/order/api", ReadOrders)
		client.POST("/order/create", CreateOrder)
		client.POST("/order/edit/", EditOrder)
		client.POST("/order/delete/:id", DeleteOrder)

		client.GET("/warehouse/api", ReadParts)
		client.POST("/warehouse/create", CreatePart)
		client.POST("/warehouse/edit/", EditPart)
		client.POST("/warehouse/delete/:id", DeletePart)
		client.POST("/warehouse/uploadImg/:id", UploadFile)
	}

	return r
}

// добавляет нужные заголовки определённым адресам, в данном случае запрашивает у браузера авторизацию для админки
func headersByRequestURI() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/public/") {
			c.Header("WWW-Authenticate", "Basic")

		}
	}
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}

//Create database connection with config
func DBConn() (db *sqlx.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "servicehelper"
	//var err error
	db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Homepage(c *gin.Context) {
	//var user User{}
	_user, _ := c.Get("user")
	user := _user.(User)

	c.HTML(http.StatusOK, "index.html",
		gin.H{
			"userLogined": user.UserName,
		},
	)

}
