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

func GetInCharges(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.InCharge{})
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
	inCharges := []models.InCharge{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.InCharge{}, fields)

	if err := db.Select(queryFields).Find(&inCharges).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(inCharges) > 0 {
		index = int(inCharges[len(inCharges)-1].ID)
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

		for _, inCharge := range inCharges {

			db.First(&inCharge.Institution, inCharge.InstitutionID)
			db.First(&inCharge.Institution.Owner, inCharge.Institution.UserID)
			db.First(&inCharge.Role, inCharge.RoleID)

			fieldMap, err := helper.FieldToMap(inCharge, fields)
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
		for _, inCharge := range inCharges {

			db.First(&inCharge.Institution, inCharge.InstitutionID)
			db.First(&inCharge.Institution.Owner, inCharge.Institution.UserID)
			db.First(&inCharge.Role, inCharge.RoleID)

			fieldMap, err := helper.FieldToMap(inCharge, fields)
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

func GetInstitutionInCharges(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.InCharge{})
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
	inCharges := []models.InCharge{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.InCharge{}, fields)

	institutionId := c.Params.ByName("id")
	if institutionId == "" {
		c.JSON(400, gin.H{"error": "Faltando id da instituição."})
		return
	}

	if err := db.Select(queryFields).Where("institution_id = ?", institutionId).Find(&inCharges).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(inCharges) > 0 {
		index = int(inCharges[len(inCharges)-1].ID)
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

		for _, inCharge := range inCharges {

			db.First(&inCharge.Institution, inCharge.InstitutionID)
			db.First(&inCharge.Institution.Owner, inCharge.Institution.UserID)
			db.First(&inCharge.Role, inCharge.RoleID)

			fieldMap, err := helper.FieldToMap(inCharge, fields)
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
		for _, inCharge := range inCharges {

			db.First(&inCharge.Institution, inCharge.InstitutionID)
			db.First(&inCharge.Institution.Owner, inCharge.Institution.UserID)
			db.First(&inCharge.Role, inCharge.RoleID)

			fieldMap, err := helper.FieldToMap(inCharge, fields)
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

func GetInCharge(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.InCharge{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	inCharge := models.InCharge{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.InCharge{}, fields)

	if err := db.Select(queryFields).First(&inCharge, id).Error; err != nil {
		content := gin.H{"error": "Encarregado com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	db.First(&inCharge.Institution, inCharge.InstitutionID)
	db.First(&inCharge.Institution.Owner, inCharge.Institution.UserID)
	db.First(&inCharge.Role, inCharge.RoleID)

	fieldMap, err := helper.FieldToMap(inCharge, fields)
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

func CreateInCharge(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	inCharge := models.InCharge{}

	if err := c.Bind(&inCharge); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if inCharge.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckInChargeMissingFields(inCharge)
	if missing != "" {
		message := "Faltando campo " + missing + " do encarregado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var user models.User
	if err = db.First(&user, inCharge.UserId).Error; err != nil {
		message := "Usuario com o id " + strconv.FormatInt(inCharge.UserId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err = db.First(&inCharge.Role, inCharge.Role.ID).Error; err != nil {
		message := "Cargo com id " + strconv.FormatInt(inCharge.Role.ID, 10) + " nao encontrado"
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err = db.First(&inCharge.Institution, inCharge.Institution.ID).Error; err != nil {
		message := "Instituição com id " + strconv.FormatInt(inCharge.Institution.ID, 10) + " não encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Create(&inCharge).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, inCharge)
}

func UpdateInCharge(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	inCharge := models.InCharge{}

	if db.First(&inCharge, id).Error != nil {
		content := gin.H{"error": "Encarregado com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&inCharge); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var institutionId int64 = inCharge.InstitutionID
	if inCharge.Institution.ID != 0 {
		institutionId = inCharge.Institution.ID
	}

	if err = db.First(&inCharge.Institution, institutionId).Error; err != nil {
		message := "Insttuicao com id " + strconv.FormatInt(institutionId, 10) + " nao encontrada."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var roleId int64 = inCharge.RoleID
	if inCharge.Role.ID != 0 {
		roleId = inCharge.Role.ID
	}

	if err = db.First(&inCharge.Role, roleId).Error; err != nil {
		message := "Cargo com id " + strconv.FormatInt(roleId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckInChargeMissingFields(inCharge)
	if missing != "" {
		message := "Faltando campo " + missing + " do encarregado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&inCharge).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	db.First(&inCharge.Institution.Owner, inCharge.Institution.UserID)

	c.JSON(200, inCharge)
}

func DeleteInCharge(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	inCharge := models.InCharge{}

	if db.First(&inCharge, id).Error != nil {
		content := gin.H{"error": "Encarregado com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&inCharge).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckInChargeMissingFields(incharge models.InCharge) string {

	if incharge.Email == "" {
		return "email"
	}

	if incharge.Name == "" {
		return "nome (name)"
	}

	if incharge.UserId == 0 {
		return "id do usuario (user_id)"
	}

	if incharge.Role.ID == 0 {
		return "id do cargo (\"role\" {\"id\": id})"
	}

	if incharge.Institution.ID == 0 {
		return "id da instituição (\"institution\": {\"id\": id})"
	}

	return ""
}

func CheckInChargeWithoutUserMissingFields(incharge models.InCharge) string {

	if incharge.Email == "" {
		return "email"
	}

	if incharge.Name == "" {
		return "nome (name)"
	}

	if incharge.Role.ID == 0 {
		return "id do cargo (\"role\" {\"id\": id})"
	}

	if incharge.Institution.ID == 0 {
		return "id da instituição (\"institution\": {\"id\": id})"
	}

	return ""
}
