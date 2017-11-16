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

func GetStudents(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Student{})
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
	students := []models.Student{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Student{}, fields)

	if err := db.Select(queryFields).Find(&students).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(students) > 0 {
		index = int(students[len(students)-1].ID)
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

		for _, student := range students {

			db.First(&student.Institution, student.InstitutionID)
			db.First(&student.Institution.Owner, student.Institution.UserID)
			db.First(&student.Responsible, student.ParentID)
			student.Institution.Owner.Password = ""

			fieldMap, err := helper.FieldToMap(student, fields)
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

		for _, student := range students {

			db.First(&student.Institution, student.InstitutionID)
			db.First(&student.Institution.Owner, student.Institution.UserID)
			db.First(&student.Responsible, student.ParentID)
			student.Institution.Owner.Password = ""

			fieldMap, err := helper.FieldToMap(student, fields)
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

func GetStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Student{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	student := models.Student{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Student{}, fields)

	if err := db.Select(queryFields).First(&student, id).Error; err != nil {
		content := gin.H{"error": "student with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	db.First(&student.Institution, student.InstitutionID)
	db.First(&student.Institution.Owner, student.Institution.UserID)
	db.First(&student.Responsible, student.ParentID)
	student.Institution.Owner.Password = ""

	fieldMap, err := helper.FieldToMap(student, fields)
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

func CreateStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	student := models.Student{}

	if err := c.Bind(&student); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	kinshipParameter := c.Params.ByName("kinship")
	if kinshipParameter == "" {
		message := "Faltando parametro de grau de parentesco(kinship.id) do criador do estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	log.Println(kinshipParameter)

	var kinshipId int64
	if kinshipId, err = strconv.ParseInt(kinshipParameter, 10, 64); err != nil {
		message := "Valor do id deve ser um inteiro (somente números). "
		c.JSON(400, gin.H{"error": message})
		return
	}

	var kinship models.Kinship
	if err = db.First(&kinship, kinshipId).Error; err != nil {
		message := "Nao foi encontrado um grau de parentesco com o id " + kinshipParameter
		c.JSON(400, gin.H{"error": message})
		return
	}

	institutionId := student.Institution.ID
	if err = db.First(&student.Institution, institutionId).Error; err != nil {
		message := "Instituicao com o id " + strconv.FormatInt(institutionId, 10) + " nao encontrada."
		c.JSON(400, gin.H{"error": message})
		return
	}

	parentId := student.Responsible.ID
	if err = db.First(&student.Responsible, parentId).Error; err != nil {
		message := "Responsavel com o id " + strconv.FormatInt(parentId, 10) + " nao encontrado(a)."
		c.JSON(400, gin.H{"error": message})
		return
	}

	db.First(&student.Institution.Owner, student.Institution.UserID)

	missing := CheckStudentMissingFields(student)
	if missing != "" {
		message := "Faltando campo de " + missing + " do estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Abre uma nova transacao
	tx := db.Begin()

	if err = tx.Create(&student).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	// Cria uma relacao entre o parente e o estudante
	parentStudent := models.ParentStudent{}
	parentStudent.Kinship = kinship
	parentStudent.Parent = student.Responsible
	parentStudent.Student = student

	if err = tx.Create(&parentStudent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	student.Institution.Owner.Password = ""

	c.JSON(201, student)
}

func UpdateStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	student := models.Student{}

	if db.First(&student, id).Error != nil {
		content := gin.H{"error": "student with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&student); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	institutionId := student.InstitutionID
	if err = db.First(&student.Institution, institutionId).Error; err != nil {
		message := "Instituicao com o id " + strconv.FormatInt(institutionId, 10) + " nao encontrada."
		c.JSON(400, gin.H{"error": message})
		return
	}

	parentId := student.ParentID
	if err = db.First(&student.Responsible, parentId).Error; err != nil {
		message := "Responsavel com o id " + strconv.FormatInt(parentId, 10) + " nao encontrado(a)."
		c.JSON(400, gin.H{"error": message})
		return
	}

	db.First(&student.Institution.Owner, student.Institution.UserID)

	missing := CheckStudentMissingFields(student)
	if missing != "" {
		message := "Faltando campo de " + missing + " do estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&student).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	student.Institution.Owner.Password = ""

	c.JSON(200, student)
}

func DeleteStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	student := models.Student{}

	if db.First(&student, id).Error != nil {
		content := gin.H{"error": "student with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&student).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckStudentMissingFields(student models.Student) string {

	if student.Name == "" {
		return "nome (name)"
	}

	if student.Responsible.ID == 0 {
		return "id do responsavel (\"responsible\": \"id\": id)"
	}

	if student.Institution.ID == 0 {
		return "id da instituicao (\"instituicao\": \"id\": id)"
	}

	return ""
}
