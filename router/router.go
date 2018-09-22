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
		api.GET("/classes_school_grades/:id", controllers.GetClassBySchoolGrade)
		api.POST("/classes", controllers.CreateClass)
		api.PUT("/classes/:id", controllers.UpdateClass)
		api.DELETE("/classes/:id", controllers.DeleteClass)

		api.GET("/disciplines", controllers.GetDisciplines)
		api.GET("/disciplines/:id", controllers.GetDiscipline)
		api.POST("/disciplines", controllers.CreateDiscipline)
		api.PUT("/disciplines/:id", controllers.UpdateDiscipline)
		api.DELETE("/disciplines/:id", controllers.DeleteDiscipline)

		api.GET("/in_charges", controllers.GetInCharges)
		api.GET("/in_charges/:id", controllers.GetInCharge)
		api.GET("/institution_in_charges/:id", controllers.GetInstitutionInCharges)
		api.POST("/in_charges", controllers.CreateInCharge)
		api.PUT("/in_charges/:id", controllers.UpdateInCharge)
		api.DELETE("/in_charges/:id", controllers.DeleteInCharge)

		api.GET("/in_charge_roles", controllers.GetInChargeRoles)
		api.GET("/in_charge_roles/:id", controllers.GetInChargeRole)
		api.POST("/in_charge_roles", controllers.CreateInChargeRole)
		api.PUT("/in_charge_roles/:id", controllers.UpdateInChargeRole)
		api.DELETE("/in_charge_roles/:id", controllers.DeleteInChargeRole)

		api.GET("/institutions", controllers.GetInstitutions)
		api.GET("/institutions/:id", controllers.GetInstitution)
		api.POST("/institutions", controllers.CreateInstitution)
		api.PUT("/institutions/:id", controllers.UpdateInstitution)
		api.DELETE("/institutions/:id", controllers.DeleteInstitution)

		api.GET("/kinships", controllers.GetKinships)
		api.GET("/kinships/:id", controllers.GetKinship)
		api.GET("/kinships_by_name/:name", controllers.GetKinshipsByName)
		api.POST("/kinships", controllers.CreateKinship)
		api.PUT("/kinships/:id", controllers.UpdateKinship)
		api.DELETE("/kinships/:id", controllers.DeleteKinship)

		api.GET("/parents", controllers.GetParents)
		api.GET("/parents/:id", controllers.GetParent)
		api.GET("/parent_user/:id", controllers.GetUserParent)
		api.POST("/parents", controllers.CreateParent)
		api.PUT("/parents/:id", controllers.UpdateParent)
		api.DELETE("/parents/:id", controllers.DeleteParent)

		api.GET("/parent_students", controllers.GetParentStudents)
		api.GET("/parent_students/:id", controllers.GetParentStudent)
		api.POST("/parent_students", controllers.CreateParentStudent)
		api.PUT("/parent_students/:id", controllers.UpdateParentStudent)
		api.DELETE("/parent_students/:id", controllers.DeleteParentStudent)

		// Registers
		api.GET("/registers", controllers.GetUserRegisters)
		api.GET("/registers/:id", controllers.GetSomeRegisters)
		api.POST("/registers", controllers.CreateRegister)
		api.POST("/registers/student", controllers.CreateRegister)
		api.PUT("/registers/:id", controllers.UpdateRegister)
		api.DELETE("/registers/:id", controllers.DeleteRegister)
		// api.GET("/registers/user/", controllers.GetUserRegisters)

		api.GET("/register_statuses", controllers.GetRegisterStatuses)
		api.GET("/register_statuses/:id", controllers.GetRegisterStatus)
		api.POST("/register_statuses", controllers.CreateRegisterStatus)
		api.PUT("/register_statuses/:id", controllers.UpdateRegisterStatus)
		api.DELETE("/register_statuses/:id", controllers.DeleteRegisterStatus)

		api.GET("/register_types", controllers.GetRegisterTypes)
		api.GET("/register_types/:id", controllers.GetRegisterType)
		api.POST("/register_types", controllers.CreateRegisterType)
		api.PUT("/register_types/:id", controllers.UpdateRegisterType)
		api.DELETE("/register_types/:id", controllers.DeleteRegisterType)

		api.GET("/school_grades", controllers.GetSchoolGrades)
		api.GET("/school_grades/:id", controllers.GetSchoolGrade)
		api.POST("/school_grades", controllers.CreateSchoolGrade)
		api.PUT("/school_grades/:id", controllers.UpdateSchoolGrade)
		api.DELETE("/school_grades/:id", controllers.DeleteSchoolGrade)

		api.GET("/students", controllers.GetStudents)
		api.GET("/students/:id", controllers.GetStudent)
		api.GET("/students/:id/class", controllers.GetStudentByClass)
		api.GET("/students_of_parent/:id", controllers.GetStudentsOfParent)
		api.POST("/students/:kinship", controllers.CreateStudent)
		api.PUT("/students/:id", controllers.UpdateStudent)
		api.DELETE("/students/:id", controllers.DeleteStudent)

		api.GET("/student_enrollments", controllers.GetStudentEnrollments)
		api.GET("/student_enrollments/:id", controllers.GetStudentEnrollment)
		api.POST("/student_enrollments", controllers.CreateStudentEnrollment)
		api.PUT("/student_enrollments/:id", controllers.UpdateStudentEnrollment)
		api.DELETE("/student_enrollments/:id", controllers.DeleteStudentEnrollment)

		api.GET("/teacher_classes", controllers.GetTeacherClasses)
		api.GET("/teacher_classes/:id", controllers.GetTeacherClass)
		api.POST("/teacher_classes", controllers.CreateTeacherClass)
		api.PUT("/teacher_classes/:id", controllers.UpdateTeacherClass)
		api.DELETE("/teacher_classes/:id", controllers.DeleteTeacherClass)

		api.GET("/users", controllers.GetUsers)
		api.GET("/users/:id", controllers.GetUser)
		api.GET("/users/:id/institutions", controllers.GetUserInstitutions)
		api.POST("/users", controllers.CreateUser)
		api.POST("/user_parent", controllers.CreateUserParent)
		api.POST("/user_incharge/:roleId/:institutionId", controllers.CreateUserInCharge)
		api.POST("/user_parent_and_student", controllers.CreateParentAndStudent)
		api.PUT("/users/:id", controllers.UpdateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)

		api.GET("/user_accesses", controllers.GetUserAccesses)
		api.GET("/user_accesses/:id", controllers.GetUserAccess)
		api.POST("/user_accesses", controllers.CreateUserAccess)
		api.PUT("/user_accesses/:id", controllers.UpdateUserAccess)
		api.DELETE("/user_accesses/:id", controllers.DeleteUserAccess)

		api.GET("/user_access_profiles", controllers.GetUserAccessProfiles)
		api.GET("/user_access_profiles/:id", controllers.GetUserAccessProfile)
		api.GET("/user_access_profiles_by_name/:name", controllers.GetUserAccessProfilesByName)
		api.POST("/user_access_profiles", controllers.CreateUserAccessProfile)
		api.PUT("/user_access_profiles/:id", controllers.UpdateUserAccessProfile)
		api.DELETE("/user_access_profiles/:id", controllers.DeleteUserAccessProfile)

		// Login
		//api.POST("/login", controllers.UserLogin)

	}

}
