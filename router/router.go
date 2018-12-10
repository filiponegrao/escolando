package router

import (
	"log"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/filiponegrao/escolando/controllers"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("secret key"),
		Timeout:         time.Hour * 24 * 7,
		MaxRefresh:      time.Hour,
		IdentityKey:     "id",
		PayloadFunc:     controllers.AuthorizationPayload,
		IdentityHandler: controllers.IdentityHandler,
		Authenticator:   controllers.UserAuthentication,
		Authorizator:    controllers.UserAuthorization,
		Unauthorized:    controllers.UserUnauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.Use(controllers.CORSMiddleware())
	r.GET("/", controllers.APIEndpoints)

	api := r.Group("")
	api.POST("/login", authMiddleware.LoginHandler)

	api.Use(authMiddleware.MiddlewareFunc())
	{

		api.GET("/classes", controllers.GetClasses)
		api.GET("/classes/:id", controllers.GetClass)
		api.GET("/classesOfSchoolGrade/:id", controllers.GetClassBySchoolGrade)
		api.GET("/classesOfInstitution/:id", controllers.GetClassByInstitution)
		api.POST("/classes", controllers.CreateClass)
		api.PUT("/classes/:id", controllers.UpdateClass)
		api.DELETE("/classes/:id", controllers.DeleteClass)

		api.GET("/disciplines", controllers.GetDisciplines)
		api.GET("/disciplines/:id", controllers.GetDiscipline)
		api.POST("/disciplines", controllers.CreateDiscipline)
		api.PUT("/disciplines/:id", controllers.UpdateDiscipline)
		api.DELETE("/disciplines/:id", controllers.DeleteDiscipline)

		api.GET("/inCharges", controllers.GetInCharges)
		api.GET("/inCharges/:id", controllers.GetInCharge)
		api.GET("/inChargesOfInstitution/:id", controllers.GetInstitutionInCharges)
		api.POST("/inCharges", controllers.CreateInCharge)
		api.POST("/inChargeUser", controllers.CreateInChargeUser)
		api.PUT("/inCharges/:id", controllers.UpdateInCharge)
		api.DELETE("/inCharges/:id", controllers.DeleteInCharge)

		api.GET("/inChargeRoles", controllers.GetInChargeRoles)
		api.GET("/inChargeRoles/:id", controllers.GetInChargeRole)
		api.POST("/inChargeRoles", controllers.CreateInChargeRole)
		api.PUT("/inChargeRoles/:id", controllers.UpdateInChargeRole)
		api.DELETE("/inChargeRoles/:id", controllers.DeleteInChargeRole)

		api.GET("/institutions", controllers.GetInstitutions)
		api.GET("/institutions/:id", controllers.GetInstitution)
		api.POST("/institutions", controllers.CreateInstitution)
		api.PUT("/institutions/:id", controllers.UpdateInstitution)
		api.DELETE("/institutions/:id", controllers.DeleteInstitution)

		api.GET("/kinships", controllers.GetKinships)
		api.GET("/kinships/:id", controllers.GetKinship)
		api.GET("/kinshipByName/:name", controllers.GetKinshipsByName)
		api.POST("/kinships", controllers.CreateKinship)
		api.PUT("/kinships/:id", controllers.UpdateKinship)
		api.DELETE("/kinships/:id", controllers.DeleteKinship)

		api.GET("/parents", controllers.GetParents)
		api.GET("/parents/:id", controllers.GetParent)
		api.GET("/parentUser", controllers.GetUserParent)
		api.GET("/parentsOfInstitution/:id", controllers.GetInstitutionParents)

		api.POST("/parents", controllers.CreateParent)
		api.PUT("/parents/:id", controllers.UpdateParent)
		api.DELETE("/parents/:id", controllers.DeleteParent)

		api.GET("/parentStudents", controllers.GetParentStudents)
		api.GET("/parentStudents/:id", controllers.GetParentStudent)
		api.POST("/parentStudents", controllers.CreateParentStudent)
		api.PUT("/parentStudents/:id", controllers.UpdateParentStudent)
		api.DELETE("/parentStudents/:id", controllers.DeleteParentStudent)

		// Registers
		api.GET("/registers", controllers.GetUserRegisters)
		api.GET("/registers/:id", controllers.GetSomeRegisters)
		api.POST("/registers", controllers.CreateRegister)
		api.POST("/registers/:parameter", controllers.CreateRegister)
		api.PUT("/registers/:id", controllers.UpdateRegister)
		api.DELETE("/registers/:id", controllers.DeleteRegister)

		api.GET("/registersChat/sent", controllers.GetSentRegistersGroupChat)
		api.DELETE("/registersChat/:id", controllers.DeleteRegisterGroupChat)

		api.GET("/registerStatus", controllers.GetRegisterStatuses)
		api.GET("/registerStatus/:id", controllers.GetRegisterStatus)
		api.POST("/registerStatus", controllers.CreateRegisterStatus)
		api.PUT("/registerStatus/:id", controllers.UpdateRegisterStatus)
		api.DELETE("/registerStatus/:id", controllers.DeleteRegisterStatus)

		api.GET("/register_types", controllers.GetRegisterTypes)
		api.GET("/register_types/:id", controllers.GetRegisterType)
		api.POST("/register_types", controllers.CreateRegisterType)
		api.PUT("/register_types/:id", controllers.UpdateRegisterType)
		api.DELETE("/register_types/:id", controllers.DeleteRegisterType)

		api.GET("/schoolGrades", controllers.GetSchoolGrades)
		api.GET("/schoolGrades/:id", controllers.GetSchoolGrade)
		api.GET("/schoolGradesOfSegment/:id", controllers.GetSchoolGradesBySegment)
		api.GET("/schoolGradesOfInstitution/:id", controllers.GetSchoolGradesByInstitution)
		api.POST("/schoolGrades", controllers.CreateSchoolGrade)
		api.PUT("/schoolGrades/:id", controllers.UpdateSchoolGrade)
		api.DELETE("/schoolGrades/:id", controllers.DeleteSchoolGrade)

		api.GET("/segments", controllers.GetAllSegments)
		api.GET("/segments/:id", controllers.GetSegments)
		api.GET("/segments_institution/:id", controllers.GetInsitutionSegments)
		api.POST("segments", controllers.CreateSegment)
		api.PUT("/segments/:id", controllers.UpdateSegment)
		api.DELETE("/segments/:id", controllers.DeleteSegment)

		api.GET("/students", controllers.GetStudents)
		api.GET("/students/:id", controllers.GetStudent)
		api.GET("/classStudents/:id", controllers.GetStudentByClass)
		api.GET("/studentsOfParent/:id", controllers.GetStudentsOfParent)
		api.GET("/studentsOfInstitution/:id", controllers.GetStudentsOfInstitution)
		api.POST("/students/:kinship", controllers.CreateStudent)
		api.PUT("/students/:id", controllers.UpdateStudent)
		api.DELETE("/students/:id", controllers.DeleteStudent)

		api.GET("/studentEnrollments", controllers.GetStudentEnrollments)
		api.GET("/studentEnrollments/:id", controllers.GetStudentEnrollment)
		api.POST("/studentEnrollments", controllers.CreateStudentEnrollment)
		api.PUT("/studentEnrollments/:id", controllers.UpdateStudentEnrollment)
		api.DELETE("/studentEnrollments/:id", controllers.DeleteStudentEnrollment)

		api.GET("/teacherClass", controllers.GetTeacherClasses)
		api.GET("/teacherClass/:id", controllers.GetTeacherClass)
		api.POST("/teacherClass", controllers.CreateTeacherClass)
		api.PUT("/teacherClass/:id", controllers.UpdateTeacherClass)
		api.DELETE("/teacherClass/:id", controllers.DeleteTeacherClass)

		api.GET("/users", controllers.GetUsers)
		api.GET("/users/:id", controllers.GetUser)
		api.GET("/users/:id/institutions", controllers.GetUserInstitutions)
		api.POST("/users", controllers.CreateUser)
		api.POST("/userParent", controllers.CreateUserParent)
		api.POST("/userParentAndStudent", controllers.CreateParentAndStudent)
		api.PUT("/users/:id", controllers.UpdateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)

		api.GET("/userAccesses", controllers.GetUserAccesses)
		api.GET("/userAccesses/:id", controllers.GetUserAccess)
		api.POST("/userAccesses", controllers.CreateUserAccess)
		api.PUT("/userAccesses/:id", controllers.UpdateUserAccess)
		api.DELETE("/userAccesses/:id", controllers.DeleteUserAccess)

		api.GET("/userAccessProfiles", controllers.GetUserAccessProfiles)
		api.GET("/userAccessProfiles/:id", controllers.GetUserAccessProfile)
		api.GET("/userAccessProfilesByName/:name", controllers.GetUserAccessProfilesByName)
		api.POST("/userAccessProfiles", controllers.CreateUserAccessProfile)
		api.PUT("/userAccessProfiles/:id", controllers.UpdateUserAccessProfile)
		api.DELETE("/userAccessProfiles/:id", controllers.DeleteUserAccessProfile)

		api.GET("/registerResponses/:id", controllers.GetRegisterResponses)
		api.POST("/registerResponses", controllers.CreateRegisterResponseForRegister)
		api.DELETE("/registerResponses/:id", controllers.DeleteRegisterResponse)

		api.POST("/password", controllers.ChangePassword)

		// Login
		//api.POST("/login", controllers.UserLogin)

	}

}
