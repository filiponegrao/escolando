package controllers

import (
	"strconv"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/models"

	"github.com/gin-gonic/gin"
)

func GetParents(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var parents models.Parent
	if err := db.Find(&parents).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, parents)
}

func GetParent(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var parent models.Parent
	if err := db.First(&parent, id).Error; err != nil {
		c.JSON(400, gin.H{"error": "Parente com o id " + id + " não encontrado."})
		return
	}
	c.JSON(200, parent)
}

func GetUserParent(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var parent models.Parent
	if err := db.Where("user_id = ?", id).First(&parent).Error; err != nil {
		c.JSON(400, gin.H{"error": "Parente com o id de usuario " + id + " não encontrado."})
		return
	}
	c.JSON(200, parent)
}

func GetInstitutionParents(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var parents []models.Parent
	if err := db.Where("institution_id = ?", id).Find(&parents).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(parents); i++ {
		db.First(&parents[i].Institution, parents[i].InstitutionID)
	}
	c.JSON(200, parents)
}

func CreateParent(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	parent := models.Parent{}

	if err := c.Bind(&parent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if parent.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}
	missing := CheckParentMissingFields(parent)
	if missing != "" {
		message := "Faltando campo " + missing + " do responsavel."
		c.JSON(400, gin.H{"error": message})
		return
	}
	var user models.User
	if err := db.First(&user, parent.UserId).Error; err != nil {
		message := "Usuario com o id " + strconv.FormatInt(parent.UserId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}
	if err := db.Create(&parent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, parent)
}

func UpdateParent(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	parent := models.Parent{}
	if db.First(&parent, id).Error != nil {
		content := gin.H{"error": "parent with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}
	if err := c.Bind(&parent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	missing := CheckParentMissingFields(parent)
	if missing != "" {
		message := "Faltando campo " + missing + " do responsavel."
		c.JSON(400, gin.H{"error": message})
		return
	}
	if err := db.Save(&parent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, parent)
}

func DeleteParent(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	parent := models.Parent{}

	if db.First(&parent, id).Error != nil {
		content := gin.H{"error": "parent with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&parent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Deletado com sucesso.")
}

func CheckParentMissingFields(parent models.Parent) string {

	if parent.Email == "" {
		return "email"
	}

	if parent.Name == "" {
		return "nome (name)"
	}

	if parent.UserId == 0 {
		return "id do usuario (user_id)"
	}

	if parent.InstitutionID == 0 {
		return "id da instituição (institutionId)"
	}

	return ""
}

func CheckParentWithoutUserMissingFields(parent models.Parent) string {

	if parent.Email == "" {
		return "email"
	}

	if parent.Name == "" {
		return "nome (name)"
	}
	if parent.Institution.ID == 0 {
		return "instituição (institution.Id)"
	}

	return ""
}
