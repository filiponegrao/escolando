package controllers

import (
	"encoding/json"
	"log"
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

			// Encontra o dono
			db.First(&institution.Owner, institution.UserID)

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
			// Encontra o dono
			db.First(&institution.Owner, institution.UserID)

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

func GetUserInstitutions(c *gin.Context) {
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

	// Recupera os ids de todas as institiocoes que o usuario possui acesso
	userId := c.Params.ByName("id")
	institutionsId := []int64{}
	rows, err := db.Table("user_accesses").Select("institution_id").Where("user_id", userId).Find(&institutionsId).Rows()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			institutionsId = append(institutionsId, id)
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	log.Println(institutionsId)

	institutions := []models.Institution{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Institution{}, fields)

	if err := db.Select(queryFields).Where(institutionsId).Find(&institutions).Error; err != nil {
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

			// Encontra o dono
			db.First(&institution.Owner, institution.UserID)

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
			// Encontra o dono
			db.First(&institution.Owner, institution.UserID)

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

	err = db.First(&institution.Owner, institution.UserID).Error
	if err != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
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

	if err = c.Bind(&institution); err != nil {
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
	err = db.First(&institution.Owner, userId).Error
	if err != nil {
		content := gin.H{"error": "Usuario com o id " + strconv.FormatInt(userId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Create(&institution).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	institution.Owner.Password = ""

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
		content := gin.H{"error": "Instituicao com o id " + id + " não encontrada."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&institution); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := institution.UserID

	err = db.First(&institution.Owner, userId).Error
	if err != nil {
		content := gin.H{"error": "Usuario com o id " + strconv.FormatInt(userId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	missing := CheckInstitutionMissingField(institution)
	if missing != "" {
		message := "Faltando campo " + missing + " da instituicao"
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&institution).Error; err != nil {
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
		content := gin.H{"error": "Instituicao com o id " + id + " não encontrada."}
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
