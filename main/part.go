package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
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

// принимаем мультипарт файл через пост, key = "file"
func UploadFile(c *gin.Context) {
	uploadId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, gin.H{"type": "uploadPartImg", "result": "notCorrectNumber"})
		return
	}

	print(uploadId)

	// проверяем айдишник на корректность
	if uploadId <= 0 {
		c.JSON(200, gin.H{"type": "uploadPartImg", "result": "notCorrectNumberLessOrZero"})
		return
	}
	file, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)

	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := c.Param("id") + extension

	absPath, _ := filepath.Abs("../ServiceHelperSRV/storage/partImg/")
	log.Println(absPath + newFileName)

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, absPath+newFileName); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"type": "uploadPartImg", "result": "unableToSave"})
		return
	}

	// File saved successfully. Return proper result
	c.JSON(200, gin.H{"type": "uploadPartImg", "result": "ok"})

}
