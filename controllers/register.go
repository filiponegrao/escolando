package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

func GetRegisters(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Register{})
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
	registers := []models.Register{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Register{}, fields)

	if err := db.Select(queryFields).Find(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(registers) > 0 {
		index = int(registers[len(registers)-1].ID)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			register, err = GetRegisterSenderInformations(register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMap, err := helper.FieldToMap(register, fields)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			register, err = GetRegisterSenderInformations(register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMap, err := helper.FieldToMap(register, fields)
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

func GetParentRegisters(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Register{})
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

	userId := c.Params.ByName("user")
	studentId := c.Params.ByName("student")

	registers := []models.Register{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

	if err := db.Where("target_id = ? AND student_id = ?", userId, studentId).Find(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(registers) > 0 {
		index = int(registers[len(registers)-1].ID)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			register, err = GetRegisterSenderInformations(register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMap, err := helper.FieldToMap(register, fields)
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

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			register, err = GetRegisterSenderInformations(register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMap, err := helper.FieldToMap(register, fields)
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

func GetRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Register{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	register := models.Register{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Register{}, fields)

	if err := db.Select(queryFields).First(&register, id).Error; err != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	register, err = GetRegisterSenderInformations(register, db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fieldMap, err := helper.FieldToMap(register, fields)
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

func CreateRegister(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	register := models.Register{}

	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if register.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckRegisterMissingFields(register, 0)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.First(&register.RegisterType, register.RegisterType.ID).Error; err != nil {
		message := "Tipo de registro com o id " + strconv.FormatInt(register.RegisterType.ID, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var targetUser models.User
	targetId := register.TargetId
	if err := db.First(&targetUser, targetId).Error; err != nil {
		message := "Usuario com id " + strconv.FormatInt(targetId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var student models.Student
	studentId := register.StudentId
	if err := db.First(&student, studentId).Error; err != nil {
		message := "Estudante com id " + strconv.FormatInt(studentId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Recupera o id do status de enviado
	if err := db.Where("name = ?", REGISTER_SENT).First(&register.Status).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&register).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	c.JSON(201, register)
}

func CreateRegisterForClass(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	register := models.Register{}

	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckRegisterMissingFields(register, 1)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Em recados para a turma, o target representa o id da turma.
	classId := register.TargetId

	if err := db.First(&register.RegisterType, register.RegisterType.ID).Error; err != nil {
		message := "Tipo de registro com o id " + strconv.FormatInt(register.RegisterType.ID, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Recupera os cadastros de alunos nessa turma
	enrollments := GetStudentEnrollmentByClass(classId, c)
	for _, enrollment := range enrollments {
		// Encontra o id de cada aluno
		studentId := enrollment.StudentID

		// Recupera todos os relacionamentos familiares desse aluno especifico
		studentParentRelations := GetParentStudentByStudent(studentId, c)
		for _, studentParentRelation := range studentParentRelations {

			// Encontra o registro do parente e salva no array
			var parent models.Parent
			if err := db.First(&parent, studentParentRelation.ParentID).Error; err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Para cada parente cria o recado
			var newRegister models.Register
			newRegister.CreatedAt = register.CreatedAt
			newRegister.RegisterType.ID = register.RegisterType.ID
			newRegister.SenderId = register.SenderId
			newRegister.TargetId = parent.UserId
			newRegister.StudentId = studentId
			newRegister.Text = register.Text
			newRegister.Title = register.Title

			if err := db.Create(&newRegister).Error; err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	}
}

func UpdateRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	register := models.Register{}

	if db.First(&register, id).Error != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var statusId int64 = register.StatusId
	if register.Status.ID != 0 {
		statusId = register.Status.ID
	}

	if err = db.First(&register.Status, statusId).Error; err != nil {
		message := "Status de registro com id " + strconv.FormatInt(statusId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	var typeId int64 = register.RegisterTypeID
	if register.RegisterType.ID != 0 {
		typeId = register.RegisterType.ID
	}

	if err = db.First(&register.RegisterType, typeId).Error; err != nil {
		message := "Tipo de registro com id " + strconv.FormatInt(typeId, 10) + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := db.Save(&register).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	c.JSON(200, register)
}

func DeleteRegister(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	register := models.Register{}

	if db.First(&register, id).Error != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&register).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func GetRegisterSenderInformations(register models.Register, db *gorm.DB) (models.Register, error) {

	register.Sender.Name = "Desconhecido"
	register.Sender.Role = "Desconhecido"

	var incharge models.InCharge
	var parent models.Parent
	var err error

	// Verifica se é um recado de funcionario
	if err = db.Where("user_id = ?", register.SenderId).Find(&incharge).Error; err == nil {

		db.First(&incharge.Role, incharge.RoleID)

		register.Sender.Name = incharge.Name
		register.Sender.Role = incharge.Role.Name

		// Se é um recado de um parente
	} else if err = db.Where("user_id = ?", register.SenderId).Find(&parent).Error; err == nil {

		register.Sender.Name = parent.Name
		register.Sender.Role = "Parente"

	} else {

		var user models.User
		db.First(&user, register.SenderId)

		register.Sender.Name = user.Name
		register.Sender.Role = "Sem cargo"

	}

	return register, nil
}

func CheckRegisterMissingFields(register models.Register, targetType int) string {

	// Target Types:
	// 0 - usuario
	// 1 - turma
	// 2 - serie

	if register.RegisterType.ID == 0 {
		return "id do tipo ('register_type':{'id': <int>})"
	}
	if register.SenderId == 0 {
		return "id do remetente ('sender_id': <int>:id)"
	}
	if register.TargetId == 0 {
		return "id do destinatario ('target_id': <int>)"
	}

	// so verifica o aluno se for um recado especifico de um aluno
	if targetType == 0 {
		if register.StudentId == 0 {
			return "id do estudante ('student_id': <int>)"
		}
	}

	return ""
}
