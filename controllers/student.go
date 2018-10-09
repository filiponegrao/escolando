package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

func GetStudentsOfParent(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	parentId := c.Params.ByName("id")
	var parentStudents []models.ParentStudent
	if err := db.Where("parent_id = ?", parentId).Find(&parentStudents).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var studentsIds []int64
	for _, parentStudent := range parentStudents {
		studentsIds = append(studentsIds, parentStudent.StudentID)
	}
	var students []models.Student
	if err := db.Find(&students, studentsIds).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(students); i++ {
		db.First(&students[i].Institution, students[i].InstitutionID)
		db.First(&students[i].Responsible, students[i].ParentID)
	}
	c.JSON(200, students)
}

func GetStudentsOfInstitution(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	institutionId := c.Params.ByName("id")
	var students []models.Student
	if err := db.Where("institution_id = ?", institutionId).Find(&students).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(students); i++ {
		db.First(&students[i].Institution, students[i].InstitutionID)
		db.First(&students[i].Responsible, students[i].ParentID)
	}
	c.JSON(200, students)
}

func GetStudent(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/students/class") {
		GetStudentByClass(c)
	} else {
		GetStudentById(c)
	}
}

func GetStudentById(c *gin.Context) {
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

func GetStudentByClass(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	classId := c.Params.ByName("id")

	var enrollments []models.StudentEnrollment
	if err := db.Where("class_id = ?", classId).Find(&enrollments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var studentIds []int64
	for _, enrollment := range enrollments {
		studentIds = append(studentIds, enrollment.StudentID)
	}

	var students []models.Student
	if err := db.Find(&students, studentIds).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, student := range students {
		db.First(&student.Institution, student.InstitutionID)
	}

	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

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

	if student.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckStudentMissingFields(student)
	if missing != "" {
		message := "Faltando campo de " + missing + " do estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	kinshipParameter := c.Params.ByName("kinship")
	if kinshipParameter == "" {
		message := "Faltando parametro de grau de parentesco(kinship.id) do criador do estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

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
	if student.Institution.ID != 0 {
		institutionId = student.Institution.ID
	}
	if err = db.First(&student.Institution, institutionId).Error; err != nil {
		message := "Instituicao com o id " + strconv.FormatInt(institutionId, 10) + " nao encontrada."
		c.JSON(400, gin.H{"error": message})
		return
	}

	parentId := student.ParentID
	if student.Responsible.ID != 0 {
		parentId = student.Responsible.ID
	}
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

func CheckStudentWithoutParentMissingFields(student models.Student) string {

	if student.Name == "" {
		return "nome (name)"
	}

	if student.Institution.ID == 0 {
		return "id da instituicao (\"instituicao\": \"id\": id)"
	}

	return ""
}
