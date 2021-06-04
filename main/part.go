package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Запчасть на складе

type Part struct {
	Id          int     `json:"id"          db:"item_id"`
	Name        string  `json:"name"        db:"name"`
	Photo       string  `json:"photo"       db:"photo"`
	SN          string  `json:"sn"          db:"sn"`
	Amount      float32 `json:"amount"      db:"amount"`
	Description string  `json:"description" db:"description"`
}

func ReadParts(c *gin.Context) {
	db := DBConn()

	parts := []Part{}

	err := db.Select(&parts, "SELECT item_id, name, photo, sn, amount, description FROM warehouse")
	if err != nil {
		c.JSON(500, gin.H{
			"type": "lastParts", "result": err.Error(),
		})
	}

	fmt.Println(parts)

	c.JSON(200, gin.H{"type": "lastParts", "parts": parts})
	defer db.Close()

}

func CreatePart(c *gin.Context) {
	db := DBConn()
	var part Part
	//fmt.Println(c.PostForm("id"))
	err := c.BindJSON(&part)
	if err != nil {
		c.JSON(500, gin.H{"type": "addPart", "result": "badPart"})
		panic(err.Error())
	}

	fmt.Println(part)

	res, err := db.NamedExec("INSERT INTO `warehouse` (`item_id`, `name`, `photo`, `sn`, `amount`, `description`) "+
		"VALUES (NULL, :name, :photo, :sn, :amount, :description);", &part)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "addPart", "result": "notCreated"})
			return
		}
		c.JSON(200, gin.H{"type": "addPart", "result": "ok"})
	}

	defer db.Close()
}

func EditPart(c *gin.Context) {
	db := DBConn()
	var part Part

	err := c.BindJSON(&part)
	if err != nil {
		c.JSON(500, gin.H{"type": "editPart", "result": "badRequest"})
		panic(err.Error())
	}

	if part.Id == 0 {
		c.JSON(500, gin.H{"type": "editPart", "result": "notValidId"})
		return
	}

	//userStatus, _ := c.Get("userStatus")
	//fmt.Println("STATUS IS ", userStatus)
	fmt.Println(part)

	res, err := db.NamedExec("UPDATE `warehouse` SET `name` = :name, `photo` = :photo, `sn` = :sn, `amount` = :amount, "+
		"`description` = :description WHERE `warehouse`.`item_id` = :item_id;", &part)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "editPart", "result": "nothingChanged"})
			return
		}
		c.JSON(200, gin.H{"type": "editPart", "result": "ok"})
	}

	defer db.Close()
}

func DeletePart(c *gin.Context) {
	db := DBConn()
	deleteId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, gin.H{"type": "deletePart", "result": "notCorrectNumber"})
		return
	}

	print(deleteId)

	// проверяем айдишник на корректность
	if deleteId <= 0 {
		c.JSON(200, gin.H{"type": "deletePart", "result": "notCorrectNumberLessOrZero"})
		return
	}

	res, err := db.Exec("DELETE FROM `warehouse` WHERE `warehouse`.`item_id` = ?", deleteId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "deletePart", "result": "nothingDeleted"})
			return
		}
		c.JSON(200, gin.H{"type": "deletePart", "result": "ok"})
	}

	defer db.Close()
}
