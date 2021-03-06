package controllers

import (
	"encoding/json"
	"net/http"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

const REGISTER_SENT = "register_sent"
const REGISTER_RECEIVED = "register_received"
const REGISTER_SEEN = "register_seen"

func GetRegisterStatuses(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.RegisterStatus{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	registerStatuses := []models.RegisterStatus{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.RegisterStatus{}, fields)

	if err := db.Select(queryFields).Find(&registerStatuses).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(registerStatuses) > 0 {
		index = int(registerStatuses[len(registerStatuses)-1].ID)
	}

	if err := parameter.SetHeaderLink(c, index); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, registerStatus := range registerStatuses {
			fieldMap, err := helper.FieldToMap(registerStatus, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if err := enc.Encode(fieldMap); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		fieldMaps := []map[string]interface{}{}

		for _, registerStatus := range registerStatuses {
			fieldMap, err := helper.FieldToMap(registerStatus, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMaps = append(fieldMaps, fieldMap)
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(200, fieldMaps)
		} else {
			c.JSON(200, fieldMaps)
		}
	}
}

func GetRegisterStatus(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.RegisterStatus{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	registerStatus := models.RegisterStatus{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.RegisterStatus{}, fields)

	if err := db.Select(queryFields).First(&registerStatus, id).Error; err != nil {
		content := gin.H{"error": "Status de registro com id" + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(registerStatus, fields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func CreateRegisterStatus(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	registerStatus := models.RegisterStatus{}

	if err := c.Bind(&registerStatus); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if registerStatus.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Create(&registerStatus).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, registerStatus)
}

func UpdateRegisterStatus(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	registerStatus := models.RegisterStatus{}

	if db.First(&registerStatus, id).Error != nil {
		content := gin.H{"error": "Status de registro com id" + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&registerStatus); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&registerStatus).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, registerStatus)
}

func DeleteRegisterStatus(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	registerStatus := models.RegisterStatus{}

	if db.First(&registerStatus, id).Error != nil {
		content := gin.H{"error": "Status de registro com id" + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&registerStatus).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func ChecRegisterStatusMissingField(status models.RegisterStatus) string {
	if status.Name == "" {
		return "nome"
	}

	return ""
}

func CheckDefaultRegisterStatus(db gorm.DB) error {

	statusArray := []string{REGISTER_SENT, REGISTER_RECEIVED, REGISTER_SEEN}
	// Para cada status default deinifindo na aplicacao:
	for _, status := range statusArray {
		var statusTemp models.RegisterStatus
		// Verifica se ja existe um registro deste status no banco:
		if err := db.Where("name = ?", status).First(&statusTemp).Error; err != nil {
			// Se nao houver, cria:
			if err2 := CreateInternalRegisterStatus(status, db); err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func CreateInternalRegisterStatus(name string, db gorm.DB) error {

	status := models.RegisterStatus{}
	status.Name = name

	if err := db.Create(&status).Error; err != nil {
		return err
	}

	return nil
}
