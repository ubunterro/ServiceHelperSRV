// Сервер для работы с приложением ServiceHelper
package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/public", "./public")

	client := r.Group("/api", basicAuth())
	{
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
	}

	return r
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
