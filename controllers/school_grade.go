package controllers

import (
	"encoding/json"
	"strconv"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

func GetSchoolGrades(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.SchoolGrade{})
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
	schoolGrades := []models.SchoolGrade{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.SchoolGrade{}, fields)

	if err := db.Select(queryFields).Find(&schoolGrades).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(schoolGrades) > 0 {
		index = int(schoolGrades[len(schoolGrades)-1].ID)
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

		for _, schoolGrade := range schoolGrades {
			fieldMap, err := helper.FieldToMap(schoolGrade, fields)
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

		for _, schoolGrade := range schoolGrades {
			fieldMap, err := helper.FieldToMap(schoolGrade, fields)
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

func GetSchoolGrade(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.SchoolGrade{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	schoolGrade := models.SchoolGrade{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.SchoolGrade{}, fields)

	if err := db.Select(queryFields).First(&schoolGrade, id).Error; err != nil {
		content := gin.H{"error": "school_grade with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(schoolGrade, fields)
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

func CreateSchoolGrade(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	schoolGrade := models.SchoolGrade{}

	if err := c.Bind(&schoolGrade); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if schoolGrade.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missingFields := CheckSchoolGradeMissingFields(schoolGrade)
	if missingFields != "" {
		message := "Faltando campo de " + missingFields + " da série."
		c.JSON(400, gin.H{"error": message})
		return
	}

	institutionId := schoolGrade.Institution.ID
	err = db.First(&schoolGrade.Institution, institutionId).Error
	if err != nil {
		content := gin.H{"error": "Instituição com o id " + strconv.FormatInt(institutionId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Create(&schoolGrade).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, schoolGrade)
}

func UpdateSchoolGrade(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	schoolGrade := models.SchoolGrade{}

	if db.First(&schoolGrade, id).Error != nil {
		content := gin.H{"error": "school_grade with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&schoolGrade); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var institutionId int64 = schoolGrade.InstitutionID
	if schoolGrade.Institution.ID != 0 {
		institutionId = schoolGrade.Institution.ID
	}

	err = db.First(&schoolGrade.Institution, institutionId).Error
	if err != nil {
		content := gin.H{"error": "Instituição com o id " + strconv.FormatInt(institutionId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	missing := CheckSchoolGradeMissingFields(schoolGrade)
	if missing != "" {
		message := "Faltando campo de " + missing + " da série."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&schoolGrade).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, schoolGrade)
}

func DeleteSchoolGrade(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	schoolGrade := models.SchoolGrade{}

	if db.First(&schoolGrade, id).Error != nil {
		content := gin.H{"error": "school_grade with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&schoolGrade).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, nil)
}

func CheckSchoolGradeMissingFields(grade models.SchoolGrade) string {

	if grade.Name == "" {
		return "nome"
	}

	if grade.Institution.ID == 0 {
		return "id da instituição (ex: 'institutuin': {'id':<int>})"
	}

	return ""
}
