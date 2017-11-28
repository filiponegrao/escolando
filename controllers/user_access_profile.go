package controllers

import (
	"encoding/json"
	"net/http"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

func GetUserAccessProfiles(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.UserAccessProfile{})
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
	userAccessProfiles := []models.UserAccessProfile{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.UserAccessProfile{}, fields)

	if err := db.Select(queryFields).Find(&userAccessProfiles).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(userAccessProfiles) > 0 {
		index = int(userAccessProfiles[len(userAccessProfiles)-1].ID)
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

		for _, userAccessProfile := range userAccessProfiles {
			fieldMap, err := helper.FieldToMap(userAccessProfile, fields)
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

		for _, userAccessProfile := range userAccessProfiles {
			fieldMap, err := helper.FieldToMap(userAccessProfile, fields)
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

func GetUserAccessProfilesByName(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.UserAccessProfile{})
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
	userAccessProfiles := []models.UserAccessProfile{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.UserAccessProfile{}, fields)

	term := c.Params.ByName("name")

	if term == "" {
		c.JSON(400, gin.H{"error": "Faltando termo a ser procurado no nome (name)"})
		return
	}

	predicate := "%" + term + "%"

	if err := db.Select(queryFields).Where("name LIKE ?", predicate).Find(&userAccessProfiles).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(userAccessProfiles) > 0 {
		index = int(userAccessProfiles[len(userAccessProfiles)-1].ID)
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

		for _, userAccessProfile := range userAccessProfiles {
			fieldMap, err := helper.FieldToMap(userAccessProfile, fields)
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

		for _, userAccessProfile := range userAccessProfiles {
			fieldMap, err := helper.FieldToMap(userAccessProfile, fields)
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

func GetUserAccessProfile(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.UserAccessProfile{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	userAccessProfile := models.UserAccessProfile{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.UserAccessProfile{}, fields)

	if err := db.Select(queryFields).First(&userAccessProfile, id).Error; err != nil {
		content := gin.H{"error": "Perfil de acesso de usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(userAccessProfile, fields)
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

func CreateUserAccessProfile(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	userAccessProfile := models.UserAccessProfile{}

	if err := c.Bind(&userAccessProfile); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userAccessProfile.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckProfileAccessMissingFields(userAccessProfile)
	if missing != "" {
		message := "Faltando campo " + missing + " do perfil de acesso"
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Create(&userAccessProfile).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, userAccessProfile)
}

func UpdateUserAccessProfile(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	userAccessProfile := models.UserAccessProfile{}

	if db.First(&userAccessProfile, id).Error != nil {
		content := gin.H{"error": "Perfil de acesso de usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&userAccessProfile); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckProfileAccessMissingFields(userAccessProfile)
	if missing != "" {
		message := "Faltando campo " + missing + " do perfil de acesso"
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&userAccessProfile).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, userAccessProfile)
}

func DeleteUserAccessProfile(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	userAccessProfile := models.UserAccessProfile{}

	if db.First(&userAccessProfile, id).Error != nil {
		content := gin.H{"error": "Perfil de acesso de usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&userAccessProfile).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckProfileAccessMissingFields(profile models.UserAccessProfile) string {
	if profile.Name == "" {
		return "nome (name)"
	}
	return ""
}
