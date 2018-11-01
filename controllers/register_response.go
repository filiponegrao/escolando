package controllers

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/models"
	"github.com/jinzhu/gorm"
)

func GetRegisterResponses(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var register models.Register
	if err := db.First(&register, id).Error; err != nil {
		c.JSON(400, gin.H{"error": "Registro não encontrado"})
		return
	}

	var responses []models.RegisterResponse
	if err := db.Where("register_id = ?", id).Find(&responses).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(responses); i++ {
		response := responses[i]

		var user models.User
		db.First(&user, response.SenderId)

		response.Sender.Name = user.Name
		response.Sender.Role = ""
		response.Sender.ImageURL = user.ProfileImageUrl
	}

	c.JSON(200, responses)
}

func CreateRegisterResponseForRegister(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))
	db := dbpkg.DBInstance(c)
	response := models.RegisterResponse{}
	if err := c.Bind(&response); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var register models.Register
	registerId := response.RegisterID
	if err := db.First(&register, registerId).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response.SenderId = userId
	tx := db.Begin()
	// Cria a resposta
	created, err := CreateRegisterResponse(response, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	register.ResponsesCount++
	if err := tx.Save(register).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, created)
}

func DeleteRegisterResponse(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")

	var response models.RegisterResponse
	if err := db.First(&response, id).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userId != response.SenderId {
		c.JSON(400, gin.H{"error": "Somente o dono de uma resposta pode excluí-la"})
		return
	}

	var register models.Register
	registerId := response.RegisterID
	if err := db.First(&register, registerId).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := db.Begin()
	if err := tx.Delete(&response).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": "Somente o dono de uma resposta pode excluí-la"})
		return
	}

	register.ResponsesCount--
	if err := tx.Save(register).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, true)
}

func CreateRegisterResponse(response models.RegisterResponse, db *gorm.DB) (models.RegisterResponse, error) {
	if response.Text == "" {
		return response, errors.New("Faltando texto da resposta (text)")
	} else if response.SenderId == 0 {
		return response, errors.New("Faltando id do remetente da resposta (response.SenderId)")
	} else if response.Register.ID == 0 {
		return response, errors.New("Faltando mensagem a ser respondida (response.Register)")
	} else {
		now := time.Now()
		response.Status.ID = 1
		response.CreatedAt = &now
		if err := db.Create(&response).Error; err != nil {
			return response, err
		} else {
			return response, nil
		}
	}
}
