package controllers

import (
	"encoding/json"
	"strconv"
	"net/http"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

func GetClasses(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Class{})
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
	classes := []models.Class{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Class{}, fields)

	if err := db.Select(queryFields).Find(&classes).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(classes) > 0 {
		index = int(classes[len(classes)-1].ID)
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

		for _, class := range classes {

			db.First(&class.SchoolGrade, class.SchoolGradeID)

			fieldMap, err := helper.FieldToMap(class, fields)
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

		for _, class := range classes {

			db.First(&class.SchoolGrade, class.SchoolGradeID)

			fieldMap, err := helper.FieldToMap(class, fields)
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

func GetClass(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Class{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	class := models.Class{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Class{}, fields)

	if err := db.Select(queryFields).First(&class, id).Error; err != nil {
		content := gin.H{"error": "class with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	db.First(&class.SchoolGrade, class.SchoolGradeID)

	fieldMap, err := helper.FieldToMap(class, fields)
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

func GetClassBySchoolGrade(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Class{})
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

	gradeId := c.Params.ByName("id")
	if gradeId == "" || gradeId == "0" {
		message := "Faltando id da série."
		c.JSON(400, gin.H{"error": message})
		return
	}

	classes := []models.Class{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

	if err := db.Where("school_grade_id = ?", gradeId).Find(&classes).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, class := range classes {

			db.First(&class.SchoolGrade, class.SchoolGradeID)

			fieldMap, err := helper.FieldToMap(class, fields)
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

		for _, class := range classes {

			db.First(&class.SchoolGrade, class.SchoolGradeID)

			fieldMap, err := helper.FieldToMap(class, fields)
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

func CreateClass(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	class := models.Class{}

	if err := c.Bind(&class); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if class.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missingField := CheckClassParameterMissing(class)
	if missingField != "" {
		message := "Faltando campo " + missingField + " da turma."
		c.JSON(400, gin.H{"error": message})
		return
	}

	gradeId := class.SchoolGrade.ID
	err = db.First(&class.SchoolGrade, gradeId).Error
	if err != nil {
		content := gin.H{"error": "Série com o id " + strconv.FormatInt(gradeId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Create(&class).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, class)
}

func UpdateClass(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	class := models.Class{}

	if db.First(&class, id).Error != nil {
		content := gin.H{"error": "class with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&class); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var gradeId int64 = class.SchoolGradeID
	if class.SchoolGrade.ID != 0 {
		gradeId = class.SchoolGrade.ID
	}

	err = db.First(&class.SchoolGrade, gradeId).Error
	if err != nil {
		content := gin.H{"error": "Série com o id " + strconv.FormatInt(gradeId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	missing := CheckClassParameterMissing(class)
	if missing != "" {

	}

	if err := db.Save(&class).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, class)
}

func DeleteClass(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	class := models.Class{}

	if db.First(&class, id).Error != nil {
		content := gin.H{"error": "class with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&class).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, nil)
}

func CheckClassParameterMissing(class models.Class) string {

	if class.Capacity == 0 {
		return "capacidade máxima da turma (capacity)"
	}

	if class.SchoolGrade.ID == 0 {
		return "id da série (ex: 'school_grade': {'id': <int>})"
	}

	if class.Name == "" {
		return "nome da turma"
	}

	return ""
}
