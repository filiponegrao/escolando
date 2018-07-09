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

func GetParentStudents(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.ParentStudent{})
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
	parentStudents := []models.ParentStudent{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.ParentStudent{}, fields)

	if err := db.Select(queryFields).Find(&parentStudents).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(parentStudents) > 0 {
		index = int(parentStudents[len(parentStudents)-1].ID)
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

		for _, parentStudent := range parentStudents {

			db.First(&parentStudent.Kinship, parentStudent.KinshipID)
			db.First(&parentStudent.Parent, parentStudent.ParentID)
			db.First(&parentStudent.Student, parentStudent.StudentID)
			db.First(&parentStudent.Student.Institution, parentStudent.Student.InstitutionID)
			db.First(&parentStudent.Student.Institution.Owner, parentStudent.Student.Institution.UserID)
			db.First(&parentStudent.Student.Responsible, parentStudent.Student.ParentID)

			fieldMap, err := helper.FieldToMap(parentStudent, fields)
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

		for _, parentStudent := range parentStudents {

			db.First(&parentStudent.Kinship, parentStudent.KinshipID)
			db.First(&parentStudent.Parent, parentStudent.ParentID)
			db.First(&parentStudent.Student, parentStudent.StudentID)
			db.First(&parentStudent.Student.Institution, parentStudent.Student.InstitutionID)
			db.First(&parentStudent.Student.Institution.Owner, parentStudent.Student.Institution.UserID)
			db.First(&parentStudent.Student.Responsible, parentStudent.Student.ParentID)

			fieldMap, err := helper.FieldToMap(parentStudent, fields)
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

func GetParentStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.ParentStudent{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	parentStudent := models.ParentStudent{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.ParentStudent{}, fields)

	if err := db.Select(queryFields).First(&parentStudent, id).Error; err != nil {
		content := gin.H{"error": "Relacao de parente com estudante com id " + id + " nao encontrada."}
		c.JSON(404, content)
		return
	}

	db.First(&parentStudent.Kinship, parentStudent.KinshipID)
	db.First(&parentStudent.Parent, parentStudent.ParentID)
	db.First(&parentStudent.Student, parentStudent.StudentID)
	db.First(&parentStudent.Student.Institution, parentStudent.Student.InstitutionID)
	db.First(&parentStudent.Student.Institution.Owner, parentStudent.Student.Institution.UserID)
	db.First(&parentStudent.Student.Responsible, parentStudent.Student.ParentID)

	fieldMap, err := helper.FieldToMap(parentStudent, fields)
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

func GetParentStudentByStudent(sutdentId int64, c *gin.Context) []models.ParentStudent {

	var result []models.ParentStudent

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.ParentStudent{})
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

	if err := db.Where("student_id = ?", sutdentId).Find(&result).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return result
	}

	return result
}

func CreateParentStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parentStudent := models.ParentStudent{}

	if err = c.Bind(&parentStudent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if parentStudent.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckParentStudentMissingFields(parentStudent)
	if missing != "" {
		message := "Faltando campo " + missing + " da relacao de parente com estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	parentId := parentStudent.Parent.ID
	if err = db.First(&parentStudent.Parent, parentId).Error; err != nil {
		message := "Parent com id " + strconv.FormatInt(parentId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	studentId := parentStudent.Student.ID
	if err = db.First(&parentStudent.Student, studentId).Error; err != nil {
		message := "Estudante com id " + strconv.FormatInt(studentId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	kinshipId := parentStudent.Kinship.ID
	if err = db.First(&parentStudent.Kinship, kinshipId).Error; err != nil {
		message := "Grau de parentesco com id " + strconv.FormatInt(kinshipId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	db.First(&parentStudent.Student.Institution, parentStudent.Student.InstitutionID)
	db.First(&parentStudent.Student.Institution.Owner, parentStudent.Student.Institution.UserID)
	db.First(&parentStudent.Student.Responsible, parentStudent.Student.ParentID)

	if err := db.Create(&parentStudent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, parentStudent)
}

/** Método responsável por criar um estudante, um partente e suas credenciais,
* e liga-los */
func CreateParentAndStudent(c *gin.Context) {

	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parentStudent := models.ParentStudent{}

	if err = c.Bind(&parentStudent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verifica a definicao de um id para um objeto parent novo.
	if parentStudent.Parent.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Verifica a definicao de um id para um objeto Student novo.
	if parentStudent.Student.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Verifica se algum campo esta faltando para a criacao de um Student
	missing := CheckStudentWithoutParentMissingFields(parentStudent.Student)
	if missing != "" {
		message := "Faltando campo " + missing + " do estudante da relacao de parente com estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	institutionId := parentStudent.Student.Institution.ID
	if err = db.First(&parentStudent.Student.Institution, institutionId).Error; err != nil {
		message := "Instituicao com id " + strconv.FormatInt(institutionId, 10) + " nao encontrada."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Verifica se algum campo esta faltando para a criacao de um Parent
	missing = CheckParentWithoutUserMissingFields(parentStudent.Parent)
	if missing != "" {
		message := "Faltando campo " + missing + " do parente da relacao de parente com estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Recupera o grau de parentesco
	kinshipId := parentStudent.Kinship.ID
	if err = db.First(&parentStudent.Kinship, kinshipId).Error; err != nil {
		message := "Grau de parentesco com id " + strconv.FormatInt(kinshipId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Define um novo usuario para o parente
	user := models.User{}
	user.Email = parentStudent.Parent.Email
	user.Name = parentStudent.Parent.Name
	user.Phone1 = parentStudent.Parent.Phone
	user.ProfileImageUrl = parentStudent.Parent.ProfileImageUrl

	// Abre uma transacao no banco
	tx := db.Begin()

	if err = tx.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	parentStudent.Parent.UserId = user.ID

	if err = tx.Create(&parentStudent.Parent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	parentStudent.Student.Responsible = parentStudent.Parent

	if err = tx.Create(&parentStudent).Error; err != nil {
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

	c.JSON(201, parentStudent)

}

func UpdateParentStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	parentStudent := models.ParentStudent{}

	if db.First(&parentStudent, id).Error != nil {
		content := gin.H{"error": "Relacao de parente com estudante com id " + id + " nao encontrada."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&parentStudent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var parentId int64 = parentStudent.ParentID
	if parentStudent.Parent.ID != 0 {
		parentId = parentStudent.Parent.ID
	}

	if err = db.First(&parentStudent.Parent, parentId).Error; err != nil {
		message := "Parente com id " + strconv.FormatInt(parentId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var studentId int64 = parentStudent.StudentID
	if parentStudent.Student.ID != 0 {
		studentId = parentStudent.Student.ID
	}

	if err = db.First(&parentStudent.Student, studentId).Error; err != nil {
		message := "Estudante com id " + strconv.FormatInt(studentId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var kinshipId int64 = parentStudent.KinshipID
	if parentStudent.Kinship.ID != 0 {
		kinshipId = parentStudent.Kinship.ID
	}

	if err = db.First(&parentStudent.Kinship, kinshipId).Error; err != nil {
		message := "Grau de parentesco com id " + strconv.FormatInt(kinshipId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckParentStudentMissingFields(parentStudent)
	if missing != "" {
		message := "Faltando campo " + missing + " da relacao de parente com estudante."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&parentStudent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.First(&parentStudent.Student.Institution, parentStudent.Student.InstitutionID)
	db.First(&parentStudent.Student.Institution.Owner, parentStudent.Student.Institution.UserID)
	db.First(&parentStudent.Student.Responsible, parentStudent.Student.ParentID)

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, parentStudent)
}

func DeleteParentStudent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	parentStudent := models.ParentStudent{}

	if db.First(&parentStudent, id).Error != nil {
		content := gin.H{"error": "Relacao de parente com estudante com id " + id + " nao encontrada."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&parentStudent).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func CheckParentStudentMissingFields(relation models.ParentStudent) string {
	if relation.Student.ID == 0 {
		return "id do estudante (\"student\": {\"id\": id})"
	}
	if relation.Parent.ID == 0 {
		return "id do parente (\"parent\": {\"id\": id})"
	}
	if relation.Kinship.ID == 0 {
		return "id do grau de parentesco (\"kinhsip\": {\"id\": id})"
	}
	return ""
}
