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

func GetClassBySchoolGrade(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var classes []models.Class
	gradeId := c.Params.ByName("id")
	if gradeId == "" {
		c.JSON(400, gin.H{"error": "Faltando id da série"})
		return
	}
	if err := db.Where("school_grade_id = ?", gradeId).Find(&classes).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(classes); i++ {
		db.First(&classes[i].InCharge, classes[i].InChargeID)
		db.First(&classes[i].InCharge.Role, classes[i].InCharge.RoleID)
		db.First(&classes[i].SchoolGrade, classes[i].SchoolGradeID)
		db.First(&classes[i].SchoolGrade.Segment, classes[i].SchoolGrade.SegmentId)
	}
	c.JSON(200, classes)
}

func GetClassByInstitution(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	institutionId := c.Params.ByName("id")
	if institutionId == "" {
		c.JSON(400, gin.H{"error": "Faltando id da instituição"})
		return
	}
	var segments []models.Segment
	if err := db.Where("institution_id = ?", institutionId).Find(&segments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var schoolGrades []models.SchoolGrade
	if err := db.Find(&schoolGrades, segments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var classes []models.Class
	if err := db.Find(&classes, schoolGrades).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(classes); i++ {
		db.First(&classes[i].InCharge, classes[i].InChargeID)
		db.First(&classes[i].InCharge.Role, classes[i].InCharge.RoleID)
		db.First(&classes[i].SchoolGrade, classes[i].SchoolGradeID)
		db.First(&classes[i].SchoolGrade.Segment, classes[i].SchoolGrade.SegmentId)
	}

	c.JSON(200, classes)
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

func CreateClass(c *gin.Context) {
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
	err := db.First(&class.SchoolGrade, gradeId).Error
	if err != nil {
		content := gin.H{"error": "Série com o id " + strconv.FormatInt(gradeId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	segmentId := class.SchoolGrade.SegmentId
	err = db.First(&class.SchoolGrade.Segment, segmentId).Error
	if err != nil {
		content := gin.H{"error": "Segmento não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Create(&class).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
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
		return "id da série (ex: 'schoolGrade': {'id': <int>})"
	}

	if class.Name == "" {
		return "nome da turma"
	}

	if class.InCharge.ID == 0 {
		return "id do responsável (ex: 'inCharge': {'id': <int>})"
	}

	return ""
}
