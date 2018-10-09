package controllers

import (
	"encoding/json"
	"strconv"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

func GetStudentEnrollments(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.StudentEnrollment{})
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
	studentEnrollments := []models.StudentEnrollment{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.StudentEnrollment{}, fields)

	if err := db.Select(queryFields).Find(&studentEnrollments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(studentEnrollments) > 0 {
		index = int(studentEnrollments[len(studentEnrollments)-1].ID)
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

		for _, studentEnrollment := range studentEnrollments {

			db.First(&studentEnrollment.Class, studentEnrollment.ClassID)
			db.First(&studentEnrollment.Student, studentEnrollment.StudentID)

			fieldMap, err := helper.FieldToMap(studentEnrollment, fields)
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

		for _, studentEnrollment := range studentEnrollments {

			db.First(&studentEnrollment.Class, studentEnrollment.ClassID)
			db.First(&studentEnrollment.Student, studentEnrollment.StudentID)

			fieldMap, err := helper.FieldToMap(studentEnrollment, fields)
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

func GetStudentEnrollment(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.StudentEnrollment{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	studentEnrollment := models.StudentEnrollment{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.StudentEnrollment{}, fields)

	if err := db.Select(queryFields).First(&studentEnrollment, id).Error; err != nil {
		content := gin.H{"error": "student_enrollment with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	db.First(&studentEnrollment.Class, studentEnrollment.ClassID)
	db.First(&studentEnrollment.Student, studentEnrollment.StudentID)

	fieldMap, err := helper.FieldToMap(studentEnrollment, fields)
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

func CreateStudentEnrollment(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	studentEnrollment := models.StudentEnrollment{}

	if err := c.Bind(&studentEnrollment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if studentEnrollment.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckStudentEnrollmentMissingField(studentEnrollment)
	if missing != "" {
		message := "Faltando campo de " + missing + " no cadastro do estudante a uma turma."
		c.JSON(400, gin.H{"error": message})
		return
	}

	studentId := studentEnrollment.Student.ID
	err = db.First(&studentEnrollment.Student, studentId).Error
	if err != nil {
		content := gin.H{"error": "Estudante com o id " + strconv.FormatInt(studentId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	classId := studentEnrollment.Class.ID
	err = db.First(&studentEnrollment.Class, classId).Error
	if err != nil {
		content := gin.H{"error": "Turma com o id " + strconv.FormatInt(classId, 10) + " não encontrada."}
		c.JSON(404, content)
		return
	}

	// Verfica a consistencia das informacoes dependentes
	if err := db.First(&studentEnrollment.Student.Institution, studentEnrollment.Student.InstitutionID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Instituição não encontrada."})
		return
	}

	if err := db.Create(&studentEnrollment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, studentEnrollment)
}

func GetStudentEnrollmentByClass(classId int64, c *gin.Context) []models.StudentEnrollment {

	var result []models.StudentEnrollment

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.StudentEnrollment{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return result
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return result
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)

	if err := db.Where("class_id = ?", classId).Find(&result).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return result
	}

	return result
}

func UpdateStudentEnrollment(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	studentEnrollment := models.StudentEnrollment{}

	if db.First(&studentEnrollment, id).Error != nil {
		content := gin.H{"error": "student_enrollment with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&studentEnrollment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var sutdentId int64 = studentEnrollment.StudentID
	if studentEnrollment.Student.ID != 0 {
		sutdentId = studentEnrollment.Student.ID
	}

	err = db.First(&studentEnrollment.Student, sutdentId).Error
	if err != nil {
		content := gin.H{"error": "Estudante com o id " + strconv.FormatInt(sutdentId, 10) + " não encontrado."}
		c.JSON(404, content)
		return
	}

	var classId int64 = studentEnrollment.ClassID
	if studentEnrollment.Class.ID != 0 {
		classId = studentEnrollment.Class.ID
	}

	err = db.First(&studentEnrollment.Class, classId).Error
	if err != nil {
		content := gin.H{"error": "Turma com o id " + strconv.FormatInt(classId, 10) + " não encontrada."}
		c.JSON(404, content)
		return
	}

	if err := db.Save(&studentEnrollment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, studentEnrollment)
}

func DeleteStudentEnrollment(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	studentEnrollment := models.StudentEnrollment{}

	if db.First(&studentEnrollment, id).Error != nil {
		content := gin.H{"error": "student_enrollment with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&studentEnrollment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, nil)
}

func CheckStudentEnrollmentMissingField(enrollment models.StudentEnrollment) string {

	if enrollment.Class.ID == 0 {
		return "id da turma ('class': {'id': <int>})"
	}

	if enrollment.Student.ID == 0 {
		return "id do aluno ('student': {'id': <int>})"
	}

	return ""
}
