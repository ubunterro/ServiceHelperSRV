package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// Order Заказ запчасти
type Order struct {
	Id          int       `json:"id"            db:"order_id"`
	Text        string    `json:"text"          db:"text"`
	Time        time.Time `json:"time"          db:"time"`
	UserOrdered int       `json:"user_ordered"  db:"user_ordered"`
	Username    string    `json:"name"          db:"name"`
}

// ReadOrders выводит JSON-структуру со списком всех заказов деталей
func ReadOrders(c *gin.Context) {
	db := DBConn()

	orders := []Order{}

	err := db.Select(&orders, "SELECT orders.order_id, orders.text, orders.time, orders.user_ordered, users.name "+
		"FROM orders INNER JOIN users ON orders.user_ordered = users.user_id;")
	if err != nil {
		c.JSON(500, gin.H{
			"type": "lastOrders", "result": err.Error(),
		})
	}

	fmt.Println(orders)

	c.JSON(200, gin.H{"type": "lastOrders", "orders": orders})
	defer db.Close()

}

func CreateOrder(c *gin.Context) {
	db := DBConn()
	var order Order
	//fmt.Println(c.PostForm("id"))
	err := c.BindJSON(&order)
	if err != nil {
		c.JSON(500, gin.H{"type": "addOrder", "result": "badOrder"})
		panic(err.Error())
	}

	fmt.Println(order)

	userId, _ := c.Get("userId")

	userIdInt, _ := strconv.Atoi(fmt.Sprintf("%v", userId))
	order.UserOrdered = userIdInt

	res, err := db.NamedExec("INSERT INTO `orders` (`order_id`, `text`, `time`, `user_ordered`) VALUES (NULL, :text, :time, :user_ordered);", &order)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "addOrder", "result": "notCreated"})
			return
		}
		c.JSON(200, gin.H{"type": "addOrder", "result": "ok"})
	}

	defer db.Close()
}

func EditOrder(c *gin.Context) {
	db := DBConn()
	var order Order
	//fmt.Println(c.PostForm("id"))
	err := c.BindJSON(&order)
	if err != nil {
		c.JSON(500, gin.H{"type": "editOrder", "result": "badRequest"})
		panic(err.Error())
	}

	if order.Id == 0 {
		c.JSON(500, gin.H{"type": "editOrder", "result": "notValidId"})
		return
	}

	//userStatus, _ := c.Get("userStatus")
	//fmt.Println("STATUS IS ", userStatus)
	fmt.Println(order)

	res, err := db.NamedExec("UPDATE `orders` SET `text` = :text, `time` = :time WHERE `orders`.`order_id` = :order_id;", &order)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "editOrder", "result": "nothingChanged"})
			return
		}
		c.JSON(200, gin.H{"type": "editOrder", "result": "ok"})
	}

	defer db.Close()
}

func DeleteOrder(c *gin.Context) {

	db := DBConn()
	deleteId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, gin.H{"type": "deleteOrder", "result": "notCorrectNumber"})
		return
	}

	print(deleteId)

	if deleteId <= 0 {
		c.JSON(200, gin.H{"type": "deleteOrder", "result": "notCorrectNumberLessOrZero"})
		return
	}

	res, err := db.Exec("DELETE FROM `orders` WHERE `orders`.`order_id` = ?;", deleteId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "deleteOrder", "result": "nothingDeleted"})
			return
		}
		c.JSON(200, gin.H{"type": "deleteOrder", "result": "ok"})
	}

	defer db.Close()
}
