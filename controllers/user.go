package controllers

import (
	"encoding/json"
	"net/http"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/tools"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.User{})
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
	users := []models.User{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.User{}, fields)

	if err := db.Select(queryFields).Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(users) > 0 {
		index = int(users[len(users)-1].ID)
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

		for _, user := range users {
			// Remove a senha por motivos de segurança
			user.Password = ""
			fieldMap, err := helper.FieldToMap(user, fields)
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

		for _, user := range users {
			// Remove a senha por motivos de segurança
			user.Password = ""
			fieldMap, err := helper.FieldToMap(user, fields)
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

func GetUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.User{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	user := models.User{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.User{}, fields)

	if err := db.Select(queryFields).First(&user, id).Error; err != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	// Remove a senha por motivos de segurança
	user.Password = ""

	fieldMap, err := helper.FieldToMap(user, fields)
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

func CreateUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckUserMissingFields(user)
	if missing != "" {
		message := "Faltando campo " + missing + " do usuario."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if user.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	user.Password = tools.EncryptTextSHA512(user.Password)

	if err := db.Create(&user).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: users.id" {
			c.JSON(400, gin.H{"error": "Ja existe um usuário com o id passado."})
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, user)
}

func CreateUserParent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckUserMissingFields(user)
	if missing != "" {
		message := "Faltando campo " + missing + " do usuario."
		c.JSON(400, gin.H{"error": message})
		return
	}

	user.Password = tools.EncryptTextSHA512(user.Password)

	tx := db.Begin()

	if err := tx.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	parent := models.Parent{}
	parent.Email = user.Email
	parent.Name = user.Name
	parent.Phone = user.Phone1
	parent.ProfileImageUrl = user.ProfileImageUrl
	parent.UserId = user.ID

	if err = tx.Create(&parent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, user)
}

func UpdateUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	user := models.User{}

	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	user.Password = ""

	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	user := models.User{}

	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckUserMissingFields(user models.User) string {

	if user.Name == "" {
		return "nome (name)"
	}

	if user.Email == "" {
		return "email"
	}

	if user.Password == "" {
		return "senha (password)"
	}

	return ""
}
