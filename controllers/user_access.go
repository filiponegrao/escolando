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

func GetUserAccesses(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.UserAccess{})
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
	userAccesses := []models.UserAccess{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.UserAccess{}, fields)

	if err := db.Select(queryFields).Find(&userAccesses).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(userAccesses) > 0 {
		index = int(userAccesses[len(userAccesses)-1].ID)
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

		for _, userAccess := range userAccesses {

			db.First(&userAccess.User, userAccess.UserID)
			db.First(&userAccess.Institution, userAccess.InstitutionID)
			db.First(&userAccess.Institution.Owner, userAccess.Institution.UserID)
			db.First(&userAccess.UserAccessProfile, userAccess.UserAccessProfileID)

			// Removendo senha dos usuarios por motivos de seguranca
			userAccess.User.Password = ""
			userAccess.Institution.Owner.Password = ""

			fieldMap, err := helper.FieldToMap(userAccess, fields)
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

		for _, userAccess := range userAccesses {

			db.First(&userAccess.User, userAccess.UserID)
			db.First(&userAccess.Institution, userAccess.InstitutionID)
			db.First(&userAccess.UserAccessProfile, userAccess.UserAccessProfileID)

			// Removendo senha dos usuarios por motivos de seguranca
			userAccess.User.Password = ""
			userAccess.Institution.Owner.Password = ""

			fieldMap, err := helper.FieldToMap(userAccess, fields)
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

func GetUserAccess(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.UserAccess{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	userAccess := models.UserAccess{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.UserAccess{}, fields)

	if err = db.Select(queryFields).First(&userAccess, id).Error; err != nil {
		content := gin.H{"error": "Acesso de usuario com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	db.First(&userAccess.User, userAccess.UserID)
	db.First(&userAccess.Institution, userAccess.InstitutionID)
	db.First(&userAccess.Institution.Owner, userAccess.Institution.UserID)
	db.First(&userAccess.UserAccessProfile, userAccess.UserAccessProfileID)

	// Removendo senha dos usuarios por motivos de seguranca
	userAccess.User.Password = ""
	userAccess.Institution.Owner.Password = ""

	fieldMap, err := helper.FieldToMap(userAccess, fields)
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

func CreateUserAccess(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	userAccess := models.UserAccess{}

	if err := c.Bind(&userAccess); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userAccess.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// AASERT: Verifica campos faltantes
	missing := CheckUserAccessMissingField(userAccess)
	if missing != "" {
		message := "Faltando campo " + missing + " do acesso."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// ASSERT: Verifica se o usuario existe de fato
	err = db.First(&userAccess.User, userAccess.User.ID).Error
	if err != nil {
		id := strconv.FormatInt(userAccess.User.ID, 10)
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	// ASSERT: Verifica se a instituicao existe de fato
	err = db.First(&userAccess.Institution, userAccess.Institution.ID).Error
	if err != nil {
		id := strconv.FormatInt(userAccess.Institution.ID, 10)
		content := gin.H{"error": "Instituicao com o id" + id + " não encontrada."}
		c.JSON(404, content)
		return
	}

	// ASSERT: Verifica se o perfil de acesso existe de fato
	err = db.First(&userAccess.UserAccessProfile, userAccess.UserAccessProfile.ID).Error
	if err != nil {
		id := strconv.FormatInt(userAccess.UserAccessProfile.ID, 10)
		content := gin.H{"error": "Perfil de acesso com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Create(&userAccess).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	userAccess.User.Password = ""
	userAccess.Institution.Owner.Password = ""

	c.JSON(201, userAccess)
}

func UpdateUserAccess(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	userAccess := models.UserAccess{}

	if db.First(&userAccess, id).Error != nil {
		content := gin.H{"error": "Acesso de usuario com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&userAccess); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// ASSERT: Verifica se o usuario existe de fato
	userId := userAccess.UserID
	if userAccess.User.ID != 0 {
		userId = userAccess.User.ID
	}
	err = db.First(&userAccess.User, userId).Error
	if err != nil {
		id := strconv.FormatInt(userId, 10)
		content := gin.H{"error": "Usuario com o id " + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	// ASSERT: Verifica se a instituicao existe de fato
	institutionId := userAccess.InstitutionID
	if userAccess.Institution.ID != 0 {
		institutionId = userAccess.Institution.ID
	}
	err = db.First(&userAccess.Institution, institutionId).Error
	if err != nil {
		id := strconv.FormatInt(institutionId, 10)
		content := gin.H{"error": "Instituicao com o id " + id + " não encontrada."}
		c.JSON(404, content)
		return
	}

	// Recupera o dono da instituicao. Nao ha necessidade de assertivas
	db.First(&userAccess.Institution.Owner, userAccess.Institution.UserID)

	// ASSERT: Verifica se o perfil de acesso existe de fato
	profileId := userAccess.UserAccessProfileID
	if userAccess.UserAccessProfile.ID != 0 {
		profileId = userAccess.UserAccessProfile.ID
	}
	err = db.First(&userAccess.UserAccessProfile, profileId).Error
	if err != nil {
		id := strconv.FormatInt(profileId, 10)
		content := gin.H{"error": "Perfil de acesso com o id " + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	// ASSERT: Verifica campos faltantes
	missing := CheckUserAccessMissingField(userAccess)
	if missing != "" {
		message := "Faltando campo " + missing + " do acesso."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&userAccess).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	userAccess.User.Password = ""
	userAccess.Institution.Owner.Password = ""

	c.JSON(200, userAccess)
}

func DeleteUserAccess(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	userAccess := models.UserAccess{}

	if db.First(&userAccess, id).Error != nil {
		content := gin.H{"error": "Acesso de usuario com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&userAccess).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckUserAccessMissingField(access models.UserAccess) string {
	if access.User.ID == 0 {
		return "id do usuario (user.id)"
	}

	if access.Institution.ID == 0 {
		return "id da instituicao (institution.id)"
	}

	if access.UserAccessProfile.ID == 0 {
		return "id do perfil de acesso (user_access_profile.id)"
	}

	return ""
}
