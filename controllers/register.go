package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

func GetRegisters(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Register{})
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
	registers := []models.Register{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Register{}, fields)

	if err := db.Select(queryFields).Find(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(registers) > 0 {
		index = int(registers[len(registers)-1].ID)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			fieldMap, err := helper.FieldToMap(register, fields)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			fieldMap, err := helper.FieldToMap(register, fields)
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

func GetParentRegisters(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Register{})
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

	userId := c.Params.ByName("user")
	studentId := c.Params.ByName("student")

	registers := []models.Register{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

	if err := db.Where("target_id = ? AND student_id = ?", userId, studentId).Find(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(registers) > 0 {
		index = int(registers[len(registers)-1].ID)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			fieldMap, err := helper.FieldToMap(register, fields)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			fieldMap, err := helper.FieldToMap(register, fields)
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

func GetRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Register{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	register := models.Register{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Register{}, fields)

	if err := db.Select(queryFields).First(&register, id).Error; err != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	fieldMap, err := helper.FieldToMap(register, fields)
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

func CreateRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	register := models.Register{}

	if err = c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if register.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckRegisterMissingFields(register)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err = db.First(&register.RegisterType, register.RegisterType.ID).Error; err != nil {
		message := "Tipo de registro com o id " + strconv.FormatInt(register.RegisterType.ID, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var targetUser models.User
	targetId := register.TargetId
	if err = db.First(&targetUser, targetId).Error; err != nil {
		message := "Usuario com id " + strconv.FormatInt(targetId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var student models.Student
	studentId := register.StudentId
	if err = db.First(&student, studentId).Error; err != nil {
		message := "Estudante com id " + strconv.FormatInt(studentId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Recupera o id do status de enviado
	if err = db.Where("name = ?", REGISTER_SENT).First(&register.Status).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&register).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	c.JSON(201, register)
}

func UpdateRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	register := models.Register{}

	if db.First(&register, id).Error != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var statusId int64 = register.StatusId
	if register.Status.ID != 0 {
		statusId = register.Status.ID
	}

	if err = db.First(&register.Status, statusId).Error; err != nil {
		message := "Status de registro com id " + strconv.FormatInt(statusId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var typeId int64 = register.RegisterTypeID
	if register.RegisterType.ID != 0 {
		typeId = register.RegisterType.ID
	}

	if err = db.First(&register.RegisterType, typeId).Error; err != nil {
		message := "Tipo de registro com id " + strconv.FormatInt(typeId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&register).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	c.JSON(200, register)
}

func DeleteRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	register := models.Register{}

	if db.First(&register, id).Error != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&register).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckRegisterMissingFields(register models.Register) string {
	if register.RegisterType.ID == 0 {
		return "id do tipo (\"register_type\":{\"id\": id})"
	}
	if register.SenderId == 0 {
		return "id do remetente (\"sender_id\":id)"
	}
	if register.TargetId == 0 {
		return "id do destinatario (\"target_id\": id)"
	}
	if register.StudentId == 0 {
		return "id do estudante (\"student_id\": id)"
	}

	return ""
}
