package db

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/filiponegrao/escolando/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/serenize/snaker"
)

func Connect() *gorm.DB {
	// SQLITE
	dir := filepath.Dir("db/database.db")
	db, err := gorm.Open("sqlite3", dir+"/database.db")

	// POSTGRES
	// db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=escolando password=postgres")
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}

	db.LogMode(false)

	if gin.IsDebugging() {
		db.LogMode(true)
	}

	if os.Getenv("AUTOMIGRATE") == "1" {
		db.AutoMigrate(
			&models.Class{},
			&models.Discipline{},
			&models.InCharge{},
			&models.InChargeRole{},
			&models.Institution{},
			&models.Kinship{},
			&models.Parent{},
			&models.ParentStudent{},
			&models.Register{},
			&models.RegisterStatus{},
			&models.RegisterType{},
			&models.SchoolGrade{},
			&models.Segment{},
			&models.Student{},
			&models.StudentEnrollment{},
			&models.TeacherClass{},
			&models.User{},
			&models.UserAccess{},
			&models.UserAccessProfile{},
			&models.RegisterResponse{},
		)
	}

	return db
}

func DBInstance(c *gin.Context) *gorm.DB {
	return c.MustGet("DB").(*gorm.DB)
}

func (self *Parameter) SetPreloads(db *gorm.DB) *gorm.DB {
	if self.Preloads == "" {
		return db
	}

	for _, preload := range strings.Split(self.Preloads, ",") {
		var a []string

		for _, s := range strings.Split(preload, ".") {
			a = append(a, snaker.SnakeToCamel(s))
		}

		db = db.Preload(strings.Join(a, "."))
	}

	return db
}
