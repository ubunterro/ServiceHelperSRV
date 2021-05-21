package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/public", "./public")

	client := r.Group("/api")
	{
		client.GET("/api", Read)
		client.POST("/category/create", Create)
		client.PATCH("/category/update/:id", Update)
		client.DELETE("/category/:id", Delete)
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
	db, err := sqlx.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type Repair struct {
	Id                  int    `json:"id"     db:"doc_id"`
	ClientName          string `json:"client" db:"client_name"`
	RepairStateId       int    `json:"status" db:"repair_state_id"`
	ProductName         string `json:"name"   db:"product_name"`
	ProductSN           string `json:"sn,"    db:"product_sn"`
	DefectDescription   string `json:"def"    db:"defect_description"`
	ReceivedWithProduct string `json:"recv"   db:"received_with_product"`
	Description         string `json:"desc"   db:"decription"`
}

func Read(c *gin.Context) {

	repairs := []Repair{}
	//repair := Repair{}

	db := DBConn()
	//rows, err := db.Query("SELECT * FROM `repairs` WHERE CategoryID = ?", c.Param("id"))
	err := db.Select(&repairs, "SELECT repairs.doc_id, clients.client_name, repairs.repair_state_id, products.product_name, products.product_sn, repairs.defect_description, repairs.received_with_product"+
		" FROM `repairs` repairs"+
		" INNER JOIN clients ON repairs.client_id = clients.client_id"+
		" INNER JOIN products ON repairs.product_id = products.product_id")
	if err != nil {
		c.JSON(500, gin.H{
			"messages": "Category not found",
		})
	}

	fmt.Println(repairs)

	c.JSON(200, gin.H{"type": "lastRepairs", "repairs": repairs})
	defer db.Close()
	//Delay close database until Read() complete
}

//Create new category API
func Create(c *gin.Context) {
	db := DBConn()

	type CreateCategory struct {
		name        string `form:"name" json:"title" binding:"required"`
		description string `form:"description" json:"body" binding:"required"`
	}

	var json CreateCategory

	if err := c.ShouldBindJSON(&json); err == nil {
		insCategory, err := db.Prepare("INSERT INTO Categories(CategoryName, Description) VALUES(?,?)")
		if err != nil {
			c.JSON(500, gin.H{
				"messages": err,
			})
		}

		insCategory.Exec(json.name, json.description)
		c.JSON(200, gin.H{
			"messages": "new category inserted",
		})

	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	defer db.Close()
}

func Update(c *gin.Context) {
	db := DBConn()
	type UpdateCategory struct {
		Title string `form:"name" json:"title" binding:"required"`
		Body  string `form:"description" json:"body" binding:"required"`
	}

	var json UpdateCategory
	if err := c.ShouldBindJSON(&json); err == nil {
		edit, err := db.Prepare("UPDATE Category SET CategoryName=?, Description=? WHERE CategoryID= " + c.Param("id"))
		if err != nil {
			panic(err.Error())
		}
		edit.Exec(json.Title, json.Body)

		c.JSON(200, gin.H{
			"messages": "category was edited",
		})
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	defer db.Close()
}

func Delete(c *gin.Context) {
	db := DBConn()

	delete, err := db.Prepare("DELETE FROM Category WHERE CategoryIDç=?")
	if err != nil {
		panic(err.Error())
	}

	delete.Exec(c.Param("id"))
	c.JSON(200, gin.H{
		"messages": "category was deleted",
	})

	defer db.Close()
}
