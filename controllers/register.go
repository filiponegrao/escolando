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

func GetSentRegistersChat(c *gin.Context) {
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

	var chats []models.Chat
	for i := 0; i < len(registers); i++ {
		register := registers[i]
		var responses []models.RegisterResponse

		if err := db.Where("register_id = ?", register.ID).Find(&responses).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var chat models.Chat
		chat.Register = register
		chat.Responses = responses

		chats = append(chats, chat)
	}

	c.JSON(200, chats)
}

func GetSentRegistersGroupChat(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))
	institutionId := c.Params.ByName("institutionId")

	var chatGroupsIds []string
	rows, err := db.Table("registers").Where("sender_id = ? AND institution_id = ?", userId, institutionId).Order("created_at desc").Select("distinct group_target_id").Rows()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var chatGroupId string
		rows.Scan(&chatGroupId)
		chatGroupsIds = append(chatGroupsIds, chatGroupId)
	}

	var chatGroups []models.ChatGroup
	for i := 0; i < len(chatGroupsIds); i++ {
		chatGroupId := chatGroupsIds[i]
		var registers []models.Register
		if err := db.Where("group_target_id = ? AND institution_id = ?", chatGroupId, institutionId).Order("created_at desc").Find(&registers).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var virtualRegister models.Register
		virtualRegister.Text = registers[0].Text
		virtualRegister.CreatedAt = registers[0].CreatedAt
		virtualRegister.RegisterType = registers[0].RegisterType
		virtualRegister.Title = registers[0].Title
		virtualRegister.GroupTargetId = chatGroupId

		// Mensagem individual
		if chatGroupId == "" {
			for _, register := range registers {
				var chatGroup models.ChatGroup
				if err := GetRegisterTargetInformations(&register, db); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				chatGroup.TargetName = register.Target.Name
				chatGroup.Register = register
				chatGroup.Targets = append(chatGroup.Targets, register.Target)
				chatGroups = append(chatGroups, chatGroup)
			}
		} else {
			// Mensagem para grupo
			var chatGroup models.ChatGroup

			targetType := strings.Split(chatGroupId, "_")[0]
			targetId := strings.Split(chatGroupId, "_")[1]
			targetName := ""

			if targetType == "SEGMENT" {
				db.Table("segments").Where("id == ?", targetId).Select("distinct(name)").Row().Scan(&targetName)
			} else if targetType == "SCHOOLGRADE" {
				db.Table("school_grades").Where("id == ?", targetId).Select("distinct(name)").Row().Scan(&targetName)
			} else if targetType == "CLASS" {
				db.Table("classes").Where("id == ?", targetId).Select("distinct(name)").Row().Scan(&targetName)
			} else if targetType == "STUDENT" {
				db.Table("students").Where("id == ?", targetId).Select("distinct(name)").Row().Scan(&targetName)
				targetName = "Parentes de " + targetName
			}

			chatGroup.Register = virtualRegister
			chatGroup.TargetName = targetName

			for _, register := range registers {
				if err := GetRegisterTargetInformations(&register, db); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				chatGroup.Targets = append(chatGroup.Targets, register.Target)
			}
			chatGroups = append(chatGroups, chatGroup)
		}
	}
	c.JSON(200, chatGroups)
}

func DeleteRegisterGroupChat(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))

	id := c.Params.ByName("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "Estrutura precisa conter o campo Registro(\"registro\")"})
		return
	}
	var registers []models.Register
	if err := db.Where("group_target_id = ?", id).Find(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Verifica se nao é o criador das mensagens
	for _, register := range registers {
		if int(register.SenderId) != userId {
			c.JSON(400, gin.H{"error": "Somente o remetente da mensagem pode deletá-la"})
			return
		}
	}
	if err := db.Delete(&registers).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
}

func GetReceivedRegisters(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))

	params := c.Request.URL.Query()
	studentId := params.Get("sutdentId")
	institutionId := params.Get("institutionId")

	if institutionId == "" {
		c.JSON(400, gin.H{"error": "Faltando id da instituicao"})
		return
	}
	var registers []models.Register

	// Recados recebidos a um aluno especifico
	if studentId != "" {
		err := db.Where("target_id = ? AND student_id = ? AND institution_id = ?", userId, studentId, institutionId).Find(&registers).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Recados recebidos, de todos os alunos.
		err := db.Where("target_id = ? AND institution_id = ?", userId, institutionId).Find(&registers).Error
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
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/sentChat") {
		GetSentRegistersChat(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/received") {
		GetReceivedRegisters(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/registers/receivedChat") {

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

	if register.SenderId == 0 {
		register.SenderId = int64(userId)
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

	// var student models.Student
	// studentId := register.StudentId
	// if err := db.First(&student, studentId).Error; err != nil {
	// 	message := "Estudante com id " + strconv.FormatInt(studentId, 10) + " nao encontrado."
	// 	c.JSON(400, gin.H{"error": message})
	// 	return
	// }

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

	uniqueID := CreateUniqueIdForGroupMessageTargets("STUDENT", studentId)
	register.GroupTargetId = uniqueID
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
	classId := register.TargetId
	// Recupera os cadastros de alunos nessa turma
	enrollments := GetStudentEnrollmentByClass(classId, c)

	uniqueID := CreateUniqueIdForGroupMessageTargets("CLASS", classId)
	register.GroupTargetId = uniqueID
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
	uniqueID := CreateUniqueIdForGroupMessageTargets("SCHOOLGRADE", schoolGradeId)
	register.GroupTargetId = uniqueID
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
	uniqueID := CreateUniqueIdForGroupMessageTargets("SEGMENT", segmentId)
	register.GroupTargetId = uniqueID

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

	newRegister.GroupTargetId = register.GroupTargetId
	newRegister.InstitutionId = register.InstitutionId

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
		register.Sender.ImageURL = incharge.ProfileImageUrl

		// Se é um recado de um parente
	} else if err = db.Where("user_id = ?", register.SenderId).Find(&parent).Error; err == nil {

		register.Sender.Name = parent.Name
		register.Sender.Role = "Parente"
		register.Sender.ImageURL = parent.ProfileImageUrl

	} else {
		var user models.User
		db.First(&user, register.SenderId)

		register.Sender.Name = user.Name
		register.Sender.Role = "Sem cargo"
		register.Sender.ImageURL = user.ProfileImageUrl

	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)
	// db.Table("register_responses").Where("register_id = ?", register.ID).Select("sum(*)").Scan(&register.ResponsesCount)

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
		register.Target.ImageURL = incharge.ProfileImageUrl

		// Se é um recado de um parente
	} else if err = db.Where("user_id = ?", register.TargetId).Find(&parent).Error; err == nil {

		register.Target.Name = parent.Name
		register.Target.Role = "Parente"
		register.Target.ImageURL = parent.ProfileImageUrl

	} else {

		var user models.User
		db.First(&user, register.TargetId)

		register.Target.Name = user.Name
		register.Target.Role = "Sem cargo"
		register.Target.ImageURL = user.ProfileImageUrl

	}

	db.First(&register.RegisterType, register.RegisterTypeID)
	db.First(&register.Status, register.StatusId)
	// db.Table("register_responses").Where("register_id = ?", register.ID).Select("sum(*)").Scan(&register.ResponsesCount)

	return nil
}

func CreateUniqueIdForGroupMessageTargets(targetType string, id int64) string {
	now := time.Now()
	uniqueID := targetType + "_" + strconv.FormatInt(id, 10) + "_" + now.String()
	uniqueID = strings.Replace(uniqueID, " ", "_", -1)
	// encoded := tools.EncryptTextSHA512(uniqueID)
	return uniqueID
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
	if register.InstitutionId == 0 {
		return "id da instituicao (institutionId: number)"
	}
	return ""
}
