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

func GetInstitutions(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Institution{})
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
	institutions := []models.Institution{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Institution{}, fields)

	if err := db.Select(queryFields).Find(&institutions).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(institutions) > 0 {
		index = int(institutions[len(institutions)-1].ID)
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

		for _, institution := range institutions {
			// Remove a senha do dono da institiocao por motivos de seguranca
			institution.Owner.Password = ""
			fieldMap, err := helper.FieldToMap(institution, fields)
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

		for _, institution := range institutions {
			// Remove a senha do dono da institiocao por motivos de seguranca
			institution.Owner.Password = ""
			fieldMap, err := helper.FieldToMap(institution, fields)
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

func GetInstitution(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Institution{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	institution := models.Institution{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Institution{}, fields)

	if err := db.Select(queryFields).First(&institution, id).Error; err != nil {
		content := gin.H{"error": "Instituicao com o id" + id + " não encontrada."}
		c.JSON(404, content)
		return
	}

	// Remove a senha do dono da institiocao por motivos de seguranca
	institution.Owner.Password = ""

	fieldMap, err := helper.FieldToMap(institution, fields)
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

func CreateInstitution(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	institution := models.Institution{}

	if err := c.Bind(&institution); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckInstitutionMissingField(institution)
	if missing != "" {
		message := "Faltando campo " + missing + " da instituicao"
		c.JSON(400, gin.H{"error": message})
		return
	}

	userId := institution.Owner.ID
	var user models.User
	err = db.First(&user, userId).Error
	if err != nil {
		content := gin.H{"error": "Usuario com o id" + strconv.FormatInt(userId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Set("gorm:save_associations", false).Create(&institution).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, institution)
}

func UpdateInstitution(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	institution := models.Institution{}

	if db.First(&institution, id).Error != nil {
		content := gin.H{"error": "Instituicao com o id" + id + " não encontrada."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&institution); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckInstitutionMissingField(institution)
	if missing != "" {
		message := "Faltando campo " + missing + " da instituicao"
		c.JSON(400, gin.H{"error": message})
		return
	}

	userId := institution.Owner.ID
	var user models.User
	err = db.First(&user, userId).Error
	if err != nil {
		content := gin.H{"error": "Usuario com o id" + strconv.FormatInt(userId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Set("gorm:save_associations", false).Save(&institution).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, institution)
}

func DeleteInstitution(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	institution := models.Institution{}

	if db.First(&institution, id).Error != nil {
		content := gin.H{"error": "Instituicao com o id" + id + " não encontrada."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&institution).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckInstitutionMissingField(institution models.Institution) string {

	if institution.Name == "" {
		return "nome (name)"
	}

	if institution.Email == "" {
		return "email"
	}

	if institution.Owner.ID == 0 {
		return "id (owner.id) do dono"
	}

	return ""
}
