package controllers

import (
	"strconv"
	"strings"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/models"
	"github.com/gin-gonic/gin"
)

func GetSegments(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "segments/institution") {
		GetInsitutionSegments(c)
	} else {
		GetSegment(c)
	}
}

func GetSegment(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Faltando id do grupo"})
		return
	}
	var segment models.Segment
	if err := db.First(&segment, id).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	db.First(&segment.Institution, segment.InstitutionID)
	c.JSON(200, segment)
}

func GetInsitutionSegments(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var segments []models.Segment
	institutionId := c.Params.ByName("id")
	if institutionId == "" {
		c.JSON(400, gin.H{"error": "Faltando id do grupo"})
		return
	}
	if err := db.Where("institution_id = ?", institutionId).Find(&segments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, segments)
}

func GetAllSegments(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var segments []models.Segment
	if err := db.Find(&segments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for index := 0; index < len(segments); index++ {
		db.First(&segments[index].Institution, segments[index].InstitutionID)
	}
	c.JSON(200, segments)
}

func CreateSegment(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	var segment models.Segment

	if err := c.Bind(&segment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckSegmentMissingField(segment)
	if missing != "" {
		message := "Faltando " + missing + " do segmento."
		c.JSON(400, gin.H{"error": message})
		return
	}

	institutionId := segment.Institution.ID
	err := db.First(&segment.Institution, institutionId).Error
	if err != nil {
		content := gin.H{"error": "Instituição com o id " + strconv.FormatInt(institutionId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Create(&segment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, segment)
}

func UpdateSegment(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Faltando id do segmento"})
		return
	}
	var segment models.Segment
	if db.First(&segment, id).Error != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}
	if err := c.Bind(&segment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := db.Save(&segment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, segment)
}

func DeleteSegment(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Faltando id do segmento"})
		return
	}
	var segment models.Segment
	if err := db.First(&segment, id).Error; err != nil {
		c.JSON(400, gin.H{"error": "Segmento com o id " + id + " não encontrado."})
		return
	}

	if err := db.Delete(&segment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "Segmento excluido com sucesso")
}
func CheckSegmentMissingField(segment models.Segment) string {
	if segment.Institution.ID == 0 {
		return "id da instituição (institutionId)"
	} else if segment.Name == "" {
		return "nome (name)"
	} else {
		return ""
	}
}
