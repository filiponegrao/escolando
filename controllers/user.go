package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	dbpkg "github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/helper"
	"github.com/filiponegrao/escolando/models"
	"github.com/filiponegrao/escolando/tools"
	"github.com/filiponegrao/escolando/version"

	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password"`
}

type newPassword struct {
	OldPassowrd     string `form:"oldPassword" json:"oldPassword"`
	NewPassword     string `form:"newPassword" json:"newPassword"`
	ConfirmPassowrd string `form:"confirmPassowrd" json:"confirmPassowrd"`
}

func GetUsers(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.User{})
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
	users := []models.User{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.User{}, fields)

	if err := db.Select(queryFields).Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(users) > 0 {
		index = int(users[len(users)-1].ID)
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

		for _, user := range users {
			// Remove a senha por motivos de segurança
			user.Password = ""
			fieldMap, err := helper.FieldToMap(user, fields)
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

		for _, user := range users {
			// Remove a senha por motivos de segurança
			user.Password = ""
			fieldMap, err := helper.FieldToMap(user, fields)
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

func GetUser(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/users/email") {
		GetUserByEmail(c)
	} else {
		GetUserById(c)
	}
}

func GetUserById(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.User{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	user := models.User{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.User{}, fields)

	if err := db.Select(queryFields).First(&user, id).Error; err != nil {
		content := gin.H{"error": "Usuario com o id " + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	// Remove a senha por motivos de segurança
	user.Password = ""

	fieldMap, err := helper.FieldToMap(user, fields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func GetUserByEmail(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	email := c.Params.ByName("email")

	user := models.User{}
	if err := db.First(&user).Where("email = ?", email).Error; err != nil {
		content := gin.H{"error": "Usuario com o email " + email + " não encontrado."}
		c.JSON(404, content)
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func CreateUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := CheckUserMissingFields(user)
	if missing != "" {
		message := "Faltando campo " + missing + " do usuario."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if user.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	user.Password = tools.EncryptTextSHA512(user.Password)

	if err := db.Create(&user).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: users.id" {
			c.JSON(400, gin.H{"error": "Ja existe um usuário com o id passado."})
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, user)
}

func CreateUserParent(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckUserMissingFields(user)
	if missing != "" {
		message := "Faltando campo " + missing + " do usuario."
		c.JSON(400, gin.H{"error": message})
		return
	}

	user.Password = tools.EncryptTextSHA512(user.Password)

	tx := db.Begin()

	if err := tx.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	parent := models.Parent{}
	parent.Email = user.Email
	parent.Name = user.Name
	parent.Phone = user.Phone1
	parent.ProfileImageUrl = user.ProfileImageUrl
	parent.UserId = user.ID

	if err = tx.Create(&parent).Error; err != nil {
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

	c.JSON(201, user)
}

func CreateUserInCharge(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	user := models.User{}

	roleId := c.Params.ByName("roleId")
	if roleId == "" {
		c.JSON(400, gin.H{"error": "Faltando id do cargo (url/user_incharge/:roleId)"})
		return
	}

	var role models.InChargeRole
	if err = db.First(&role, roleId).Error; err != nil {
		message := "Cargo com o id " + roleId + " não encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	institutionId := c.Params.ByName("institutionId")
	if institutionId == "" {
		c.JSON(400, gin.H{"error": "Faltando id da instituição (url/user_incharge/:roleId/:institutionId)"})
		return
	}

	var institution models.Institution
	if err = db.First(&institution, institutionId).Error; err != nil {
		message := "Instituicao com o id " + roleId + " não encontrada."
		c.JSON(400, gin.H{"error": message})
		return
	}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.ID != 0 {
		message := "Nao é permitida a escolha de um id para um novo objeto."
		c.JSON(400, gin.H{"error": message})
		return
	}

	missing := CheckUserMissingFields(user)
	if missing != "" {
		message := "Faltando campo " + missing + " do usuario."
		c.JSON(400, gin.H{"error": message})
		return
	}

	user.Password = tools.EncryptTextSHA512(user.Password)

	tx := db.Begin()

	if err := tx.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	incharge := models.InCharge{}
	incharge.Email = user.Email
	incharge.Name = user.Name
	incharge.Phone = user.Phone1
	incharge.ProfileImageUrl = user.ProfileImageUrl
	incharge.UserId = user.ID
	incharge.Institution = institution
	incharge.Role = role

	if err = tx.Create(&incharge).Error; err != nil {
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

	c.JSON(201, user)
}

func UpdateUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	user := models.User{}

	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	user.Password = ""

	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	user := models.User{}

	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "Usuario com o id" + id + " não encontrado."}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func Login(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" {
		message := "Faltando email"
		c.JSON(400, gin.H{"error": message})
		return
	}

	if password == "" {
		message := "Faltando senha (password)"
		c.JSON(400, gin.H{"error": message})
		return
	}

	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		message := "Usuario com email " + email + " nao encontrado."
		c.JSON(400, gin.H{"error": message})
		return
	}

	encPassword := tools.EncryptTextSHA512(password)

	if encPassword != user.Password {
		message := "Senha incorreta"
		c.JSON(400, gin.H{"error": message})
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func ChangePassword(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var newPassword newPassword
	if err := c.Bind(&newPassword); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	claims := jwt.ExtractClaims(c)
	userId := int64(claims["id"].(float64))
	var user models.User
	if err := db.First(&user, userId).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	encPassword := tools.EncryptTextSHA512(user.Password)
	// Vrifica corretude de senha recebida
	if newPassword.OldPassowrd != encPassword {
		message := "Senha atual incorreta."
		c.JSON(400, gin.H{"error": message})
		return
	}
	// Verifica se nova senha foi escrita corretamente
	if newPassword.NewPassword != newPassword.ConfirmPassowrd {
		message := "Nova senha nao confere."
		c.JSON(400, gin.H{"error": message})
		return
	}
	newPasswordEnc := tools.EncryptTextSHA512(newPassword.NewPassword)
	user.Password = newPasswordEnc
	db.Save(user)

	c.JSON(200, "Senha atualizada com sucesso")
}

func UserAuthentication(c *gin.Context) (interface{}, error) {

	var loginVals login

	if err := c.Bind(&loginVals); err != nil {
		return nil, err
	}

	email := loginVals.Username
	password := loginVals.Password

	db := dbpkg.DBInstance(c)

	if email == "" {
		message := "Faltando email"
		c.JSON(400, gin.H{"error": message})
		return nil, errors.New(message)
	}

	// if password == "" {
	// 	message := "Faltando senha (password)"
	// 	c.JSON(400, gin.H{"error": message})
	// 	return nil, errors.New(message)
	// }

	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		//message := "Usuario com email " + email + " nao encontrado."
		//c.JSON(400, gin.H{"error": message})
		return nil, err
	}

	// Verifica se é usuario com primeiro acesso
	if password == "" {

	} else {
		encPassword := tools.EncryptTextSHA512(password)

		if encPassword != user.Password {
			//message := "Senha incorreta"
			//c.JSON(400, gin.H{"error": message})
			return nil, errors.New("Senha incorreta")
		}
	}

	user.Password = ""

	return &user, nil
}

func UserAuthorization(user interface{}, c *gin.Context) bool {
	return true
}

// Falha na autênticação
func UserUnauthorized(c *gin.Context, code int, message string) {
	err := ""
	if strings.Contains(message, "missing") {
		err = "Faltando email ou senha"
	} else if strings.Contains(message, "incorrect") {
		err = "Email ou senha incorreta"
	} else if strings.Contains(message, "cookie token is empty") {
		err = "Faltando HEADER de autenticação!"
	} else {
		err = message
	}
	c.JSON(code, gin.H{"error": err})
}

// func AuthorizationPayload(data interface{}) jwt.MapClaims {
// 	m := make(map[string]interface{})
// 	if v, ok := data.(*models.User); ok {
// 		m["user_id"] = v.ID
// 	}
// 	return m
// }

func AuthorizationPayload(data interface{}) jwt.MapClaims {
	if user, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"id": user.ID,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.User{
		ID: int64(claims["id"].(float64)),
	}
}

func CheckUserMissingFields(user models.User) string {

	if user.Name == "" {
		return "nome (name)"
	}

	if user.Email == "" {
		return "email"
	}

	if user.Password == "" {
		return "senha (password)"
	}

	return ""
}
