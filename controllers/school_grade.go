package controllers

import (
	"log"
	"strconv"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"

	"github.com/gin-gonic/gin"
)

func GetSchoolGrades(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var grades []models.SchoolGrade
	if err := db.Find(&grades).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(grades); i++ {
		db.First(&grades[i].Segment, grades[i].SegmentId)
		db.First(&grades[i].Segment.Institution, grades[i].Segment.InstitutionID)

	}
	c.JSON(200, grades)
}

func GetSchoolGrade(c *gin.Context) {
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
		content := gin.H{"error": "Série com id " + id + " nao encontrada"}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(schoolGrade, fields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func CreateSchoolGrade(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	schoolGrade := models.SchoolGrade{}

	if err := c.Bind(&schoolGrade); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Println("Serie a ser criada: ", schoolGrade)

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
	segmentId := schoolGrade.Segment.ID
	err := db.First(&schoolGrade.Segment, segmentId).Error
	if err != nil {
		content := gin.H{"error": "Segmento com o id " + strconv.FormatInt(segmentId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	institutionId := schoolGrade.Segment.InstitutionID
	err = db.First(&schoolGrade.Segment.Institution, institutionId).Error
	if err != nil {
		content := gin.H{"error": "Instituião não encontrada."}
		c.JSON(404, content)
		return
	}

	log.Println("Serie a ser criada com informacoes preenchidas: ", schoolGrade)

	if err := db.Create(&schoolGrade).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, schoolGrade)
}

func UpdateSchoolGrade(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	schoolGrade := models.SchoolGrade{}

	if db.First(&schoolGrade, id).Error != nil {
		content := gin.H{"error": "Série com o id " + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&schoolGrade); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var segmentId int64 = schoolGrade.SegmentId
	if schoolGrade.Segment.ID != 0 {
		segmentId = schoolGrade.Segment.ID
	}

	err := db.First(&schoolGrade.Segment, segmentId).Error
	if err != nil {
		content := gin.H{"error": "Segmento com o id " + strconv.FormatInt(segmentId, 10) + " não encontrado."}
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

	c.JSON(200, schoolGrade)
}

func DeleteSchoolGrade(c *gin.Context) {
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
	c.JSON(200, nil)
}

func GetSchoolGradesBySegment(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var schoolGrades []models.SchoolGrade
	segmentId := c.Params.ByName("id")
	if segmentId == "" {
		c.JSON(400, gin.H{"error": "Faltando id do segmento"})
		return
	}
	if err := db.Where("segment_id = ?", segmentId).Find(&schoolGrades).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(schoolGrades); i++ {
		db.First(&schoolGrades[i].Segment, schoolGrades[i].SegmentId)
		db.First(&schoolGrades[i].Segment.Institution, schoolGrades[i].Segment.InstitutionID)
	}
	c.JSON(200, schoolGrades)
}

func GetSchoolGradesByInstitution(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var segments []models.Segment
	institutionId := c.Params.ByName("id")
	if institutionId == "" {
		c.JSON(400, gin.H{"error": "Faltando id do segmento"})
		return
	}
	if err := db.Where("institution_id = ?", institutionId).Find(&segments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var schoolGrades []models.SchoolGrade
	if err := db.Find(&schoolGrades, segments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(schoolGrades); i++ {
		db.First(&schoolGrades[i].Segment, schoolGrades[i].SegmentId)
		db.First(&schoolGrades[i].Segment.Institution, schoolGrades[i].Segment.InstitutionID)
	}

	c.JSON(200, schoolGrades)
}

func CheckSchoolGradeMissingFields(grade models.SchoolGrade) string {
	if grade.Name == "" {
		return "nome"
	}
	if grade.Segment.ID == 0 {
		return "id do segmento (ex: 'segment': {'id':<int>})"
	}
	return ""
}
