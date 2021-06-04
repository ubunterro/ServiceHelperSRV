package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Repair struct {
	Id                  int     `json:"id"     db:"doc_id"`
	ClientName          string  `json:"client" db:"client_name"`
	RepairStateId       int     `json:"status" db:"repair_state_id"`
	ProductName         string  `json:"name"   db:"product_name"`
	ProductSN           string  `json:"sn"     db:"product_sn"`
	DefectDescription   string  `json:"def"    db:"defect_description"`
	ReceivedWithProduct string  `json:"recv"   db:"received_with_product"`
	Description         string  `json:"desc"   db:"description"`
	Amount              float32 `json:"amount" db:"amount"`
}

func Read(c *gin.Context) {

	repairs := []Repair{}
	//repair := Repair{}

	db := DBConn()
	//rows, err := db.Query("SELECT * FROM `repairs` WHERE CategoryID = ?", c.Param("id"))
	err := db.Select(&repairs, "SELECT repairs.doc_id, clients.client_name, repairs.repair_state_id, products.product_name, products.product_sn, repairs.defect_description, repairs.received_with_product"+
		" FROM `repairs` repairs"+
		" INNER JOIN clients ON repairs.client_id = clients.client_id"+
		" INNER JOIN products ON repairs.product_id = products.product_id"+
		" WHERE repairs.is_deleted = 0"+
		" ORDER BY repairs.doc_id")
	if err != nil {
		c.JSON(500, gin.H{
			"type": "lastRepairs", "result": err.Error(),
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
	var repair Repair
	//fmt.Println(c.PostForm("id"))
	err := c.BindJSON(&repair)
	if err != nil {
		c.JSON(500, gin.H{"type": "addRepair", "result": "badRequest"})
		panic(err.Error())
	}

	fmt.Println(repair)
	err = nil
	//_, err = db.NamedExec("INSERT INTO `repairs` (`doc_id`, `client_id`, `repair_state_id`, `product_id`, `defect_description`, `received_with_product`, `amount`, `description`)" +
	//							 " VALUES (NULL, :ClientName, :repair_state_id, :defect_description, :received_with_product, :amount, :description)", &repair)
	tx := db.MustBegin()
	// language=SQL
	tx.NamedExec("INSERT INTO clients "+
		"SET client_name = :client_name "+
		"ON DUPLICATE KEY UPDATE client_id = client_id; ", &repair)
	tx.NamedExec("SELECT client_id INTO @client_id FROM clients WHERE client_name = :client_name;", &repair)

	tx.NamedExec("INSERT INTO products "+
		"SET product_name = :product_name, product_sn = :product_sn, description = :description ON DUPLICATE KEY UPDATE product_id = product_id;", &repair)
	tx.NamedExec("SELECT product_id INTO @product_id FROM products WHERE product_name = :product_name AND product_sn = :product_sn;", &repair)

	tx.NamedExec("INSERT INTO `repairs` (`doc_id`, `client_id`, `repair_state_id`, `product_id`, `defect_description`, `received_with_product`, `amount`, `description`) VALUES"+
		" (NULL, @client_id, :repair_state_id, @product_id, :defect_description, :received_with_product, :amount, :description);", &repair)
	tx.Commit()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	} else {
		//count, _ := res.RowsAffected()
		//if count == 0 {
		//	c.JSON(500, gin.H{"type" : "addRepair", "result" : "notFound"})
		//	return
		//}
		c.JSON(200, gin.H{"type": "addRepair", "result": "ok"})
	}

	defer db.Close()

}

func Edit(c *gin.Context) {
	// UPDATE repairs SET repair_state_id = 1, defect_description = '', received_with_product = '', amount = 0, description = 1 WHERE doc_id = 1;
	db := DBConn()
	var repair Repair
	//fmt.Println(c.PostForm("id"))
	err := c.BindJSON(&repair)
	if err != nil {
		c.JSON(500, gin.H{"type": "editRepair", "result": "badRequest"})
		panic(err.Error())
	}

	if repair.Id == 0 {
		c.JSON(500, gin.H{"type": "editRepair", "result": "notValidId"})
		return
	}

	userStatus, _ := c.Get("userStatus")
	fmt.Println("STATUS IS ", userStatus)
	fmt.Println(repair)
	err = nil

	res, err := db.NamedExec("UPDATE repairs SET repair_state_id = :repair_state_id, defect_description = :defect_description, "+
		"received_with_product = :received_with_product, amount = :amount, description = :description "+
		"WHERE doc_id = :doc_id;", &repair)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "editRepair", "result": "nothingChanged"})
			return
		}
		c.JSON(200, gin.H{"type": "editRepair", "result": "ok"})
	}

	defer db.Close()
}

func Delete(c *gin.Context) {

	db := DBConn()
	deleteId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, gin.H{"type": "delete", "result": "notCorrectNumber"})
		return
	}

	print(deleteId)

	if deleteId <= 0 {
		c.JSON(200, gin.H{"type": "delete", "result": "notCorrectNumberLessOrZero"})
		return
	}

	res, err := db.Exec("UPDATE repairs SET is_deleted = 1 WHERE doc_id = ?;", deleteId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		count, _ := res.RowsAffected()
		if count == 0 {
			c.JSON(500, gin.H{"type": "deleteRepair", "result": "nothingDeleted"})
			return
		}
		c.JSON(200, gin.H{"type": "deleteRepair", "result": "ok"})
	}

	defer db.Close()
}
