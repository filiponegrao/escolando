package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/version"
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

func GetRegisters(c *gin.Context) {
	println("PRIMEIRO PASSO:")
	if strings.HasPrefix(c.Request.RequestURI, "/registers/user") {
		GetUserRegisters(c)
	} else {
		GetAllRegisters(c)
	}
}

func GetAllRegisters(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	registers := []models.Register{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Register{}, fields)

	if err := db.Select(queryFields).Find(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, register := range registers {

			var err error

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			err = GetRegisterSenderInformations(&register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err = GetRegisterTargetInformations(&register, db)
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

			var err error

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			err = GetRegisterSenderInformations(&register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err = GetRegisterTargetInformations(&register, db)
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

func GetUserRegisters(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	claims := jwt.ExtractClaims(c)
	userId := int(claims["user_id"].(float64))

	studentId := c.Params.ByName("student")

	registers := []models.Register{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))

	// Verifica se é um recado para um aluno especifico
	if studentId == "" {
		if err := db.Where("target_id = ? OR sender_id = ?", userId, userId).Find(&registers).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Caso nao seja para um aluno especifico
		response := db.Where("(target_id = ? OR sender_id = ?) AND (student_id = ?)", userId, userId, studentId).Find(&registers)
		err := response.Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, register := range registers {

			db.First(&register.RegisterType, register.RegisterTypeID)
			db.First(&register.Status, register.StatusId)

			err := GetRegisterSenderInformations(&register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err = GetRegisterTargetInformations(&register, db)
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

			err := GetRegisterSenderInformations(&register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err = GetRegisterTargetInformations(&register, db)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
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

func GetSentRegisters(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))
	studentId := c.Params.ByName("studentId")
	var registers []models.Register

	// Recados enviados a um aluno especifico
	if studentId != "" {
		err := db.Where("sender_id = ? AND student_id = ?", userId, studentId).Find(&registers).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		err := db.Where("sender_id = ?", userId).Find(&registers).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	registersInfo := make([]models.Register, 0)

	for _, register := range registers {
		err := GetRegisterTargetInformations(&register, db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		} else {
			registersInfo = append(registersInfo, register)
		}
	}

	c.JSON(200, registersInfo)
}

func GetReceivedRegisters(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))
	studentId := c.Params.ByName("studentId")
	var registers []models.Register

	// Recados enviados a um aluno especifico
	if studentId != "" {
		err := db.Where("target_id = ? AND student_id = ?", userId, studentId).Find(&registers).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		err := db.Where("target_id = ?", userId).Find(&registers).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	registersInfo := make([]models.Register, 0)

	for _, register := range registers {
		err := GetRegisterSenderInformations(&register, db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		} else {
			registersInfo = append(registersInfo, register)
		}
	}

	c.JSON(200, registersInfo)
}

func GetSomeRegisters(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/registers/sent") {
		GetSentRegisters(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/received") {
		GetReceivedRegisters(c)
	} else {
		GetRegister(c)
	}
}

func GetRegister(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	register := models.Register{}
	id := c.Params.ByName("id")

	if err := db.First(&register, id).Error; err != nil {
		content := gin.H{"error": "Registro com id " + id + " nao encontrado."}
		c.JSON(404, content)
		return
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	err := GetRegisterSenderInformations(&register, db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = GetRegisterTargetInformations(&register, db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, register)
}

func CreateRegister(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/registers/segment") {
		CreateRegisterForSegment(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/grade") {
		CreateRegisterForSchoolGrade(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/class") {
		CreateRegisterForClass(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/student") {
		CreateRegisterForStudent(c)
	} else {
		CreateSingleRegister(c)
	}
}

func CreateSingleRegister(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))

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

	register.SenderId = int64(userId)
	println("CRIA RECADO")
	println(register.TargetId)

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

func CreateRegisterForStudent(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))

	var register models.Register
	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	register.SenderId = userId

	missing := CheckRegisterMissingFields(register, 1)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	studentId := register.TargetId

	var relations []models.ParentStudent

	if err := db.Where("student_id = ?", studentId).Find(&relations).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, relation := range relations {

		var parent models.Parent

		if err := db.First(&parent, relation.ParentID).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		register.StudentId = studentId
		register.TargetId = parent.UserId

		if err := SaveRegister(db, register, userId); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}
}

func CreateRegisterForClass(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))

	register := models.Register{}
	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	register.SenderId = userId
	missing := CheckRegisterMissingFields(register, 1)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Em recados para a turma, o target representa o id da turma.
	classId := register.TargetId
	// if err := db.First(&register.RegisterType, register.RegisterType.ID).Error; err != nil {
	// 	message := "Tipo de registro com o id " + strconv.FormatInt(register.RegisterType.ID, 10) + " nao encontrado."
	// 	c.JSON(400, gin.H{"error": message})
	// 	return
	// }
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

			register.StudentId = studentId
			register.TargetId = parent.UserId

			if err := SaveRegister(db, register, userId); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	}
}

func CreateRegisterForSchoolGrade(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))
	register := models.Register{}
	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	register.SenderId = userId
	missing := CheckRegisterMissingFields(register, 2)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}
	schoolGradeId := register.TargetId
	// Encontra todas as turmas desta série
	var classes []models.Class
	if err := db.Where("school_grade_id = ?", schoolGradeId).Find(&classes).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for _, class := range classes {
		var enrollments []models.StudentEnrollment
		if err := db.Where("class_id = ?", class.ID).Find(&enrollments).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		for _, enrollment := range enrollments {
			var relations []models.ParentStudent
			if err := db.Where("student_id = ?", enrollment.StudentID).Find(&relations).Error; err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			for _, relation := range relations {
				var parent models.Parent
				if err := db.First(&parent, relation.ParentID).Error; err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				register.StudentId = enrollment.StudentID
				register.TargetId = parent.UserId
				if err := SaveRegister(db, register, userId); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
			}
		}
	}
}

func CreateRegisterForSegment(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))
	register := models.Register{}
	if err := c.Bind(&register); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	register.SenderId = userId
	missing := CheckRegisterMissingFields(register, 2)
	if missing != "" {
		message := "Faltando campo " + missing + " do recado."
		c.JSON(400, gin.H{"error": message})
		return
	}
	segmentId := register.TargetId
	// Encontra todas as series desse segmento
	var grades []models.SchoolGrade
	if err := db.Where("segment_id = ?", segmentId).Find(&grades).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for _, grade := range grades {
		// Encontra todas as turmas de cada segmento
		var classes []models.Class
		if err := db.Where("school_grade_id = ?", grade.ID).Find(&classes).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		for _, class := range classes {
			var enrollments []models.StudentEnrollment
			if err := db.Where("class_id = ?", class.ID).Find(&enrollments).Error; err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			for _, enrollment := range enrollments {
				var relations []models.ParentStudent
				if err := db.Where("student_id = ?", enrollment.StudentID).Find(&relations).Error; err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				for _, relation := range relations {
					var parent models.Parent
					if err := db.First(&parent, relation.ParentID).Error; err != nil {
						c.JSON(400, gin.H{"error": err.Error()})
						return
					}
					register.StudentId = enrollment.StudentID
					register.TargetId = parent.UserId
					if err := SaveRegister(db, register, userId); err != nil {
						c.JSON(400, gin.H{"error": err.Error()})
						return
					}
				}
			}
		}
	}
}

func SaveRegister(db *gorm.DB, register models.Register, sender int64) error {
	var newRegister models.Register

	date := time.Now()
	newRegister.CreatedAt = &date
	newRegister.UpdatedAt = &date
	newRegister.StatusId = 1
	newRegister.RegisterTypeID = 1
	newRegister.Text = register.Text
	newRegister.Title = register.Title
	newRegister.StudentId = register.StudentId

	newRegister.SenderId = sender
	newRegister.TargetId = register.TargetId

	if err := db.Create(&newRegister).Error; err != nil {
		return err
	} else {
		return nil
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

func GetRegisterSenderInformations(register *models.Register, db *gorm.DB) error {

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

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	return nil
}

func GetRegisterTargetInformations(register *models.Register, db *gorm.DB) error {

	register.Target.Name = "Desconhecido"
	register.Target.Role = "Desconhecido"

	var incharge models.InCharge
	var parent models.Parent
	var err error

	// Verifica se é um recado de funcionario
	if err = db.Where("user_id = ?", register.TargetId).Find(&incharge).Error; err == nil {

		db.First(&incharge.Role, incharge.RoleID)

		register.Target.Name = incharge.Name
		register.Target.Role = incharge.Role.Name

		// Se é um recado de um parente
	} else if err = db.Where("user_id = ?", register.TargetId).Find(&parent).Error; err == nil {

		register.Target.Name = parent.Name
		register.Target.Role = "Parente"

	} else {

		var user models.User
		db.First(&user, register.TargetId)

		register.Target.Name = user.Name
		register.Target.Role = "Sem cargo"
	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)

	return nil
}

func CheckRegisterMissingFields(register models.Register, targetType int) string {

	// Target Types:
	// 0 - usuario
	// 1 - turma
	// 2 - serie
	// 3 - parentes

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
	if targetType == 0 || targetType == 3 {
		if register.StudentId == 0 {
			return "id do estudante ('student_id': <int>)"
		}
	}

	return ""
}
