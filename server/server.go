package server

import (
	"github.com/filiponegrao/escolando/controllers"
	"github.com/filiponegrao/escolando/middleware"
	"github.com/filiponegrao/escolando/router"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.SetDBtoContext(db))
	router.Initialize(r)

	return r
}

func InitConfigurations(db gorm.DB) error {

	var err error

	if err = controllers.CheckDefaultRegisterStatus(db); err != nil {
		return err
	}

	return nil
}
