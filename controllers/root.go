package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", reqScheme, reqHost)

	resources := map[string]string{
		"classes_url":                   baseURL + "/classes",
		"class_url":                     baseURL + "/classes/{id}",
		"disciplines_url":               baseURL + "/disciplines",
		"discipline_url":                baseURL + "/disciplines/{id}",
		"in_charges_url":                baseURL + "/in_charges",
		"in_charge_url":                 baseURL + "/in_charges/{id}",
		"in_charge_roles_url":           baseURL + "/in_charge_roles",
		"in_charge_role_url":            baseURL + "/in_charge_roles/{id}",
		"institutions_url":              baseURL + "/institutions",
		"institution_url":               baseURL + "/institutions/{id}",
		"kinships_url":                  baseURL + "/kinships",
		"kinship_url":                   baseURL + "/kinships/{id}",
		"parents_url":                   baseURL + "/parents",
		"parent_url":                    baseURL + "/parents/{id}",
		"parent_students_url":           baseURL + "/parent_students",
		"parent_student_url":            baseURL + "/parent_students/{id}",
		"registers_url":                 baseURL + "/registers",
		"register_url":                  baseURL + "/registers/{id}",
		"register_current_statuses_url": baseURL + "/register_current_statuses",
		"register_current_status_url":   baseURL + "/register_current_statuses/{id}",
		"register_statuses_url":         baseURL + "/register_statuses",
		"register_status_url":           baseURL + "/register_statuses/{id}",
		"register_types_url":            baseURL + "/register_types",
		"register_type_url":             baseURL + "/register_types/{id}",
		"school_grades_url":             baseURL + "/school_grades",
		"school_grade_url":              baseURL + "/school_grades/{id}",
		"segments_url":                  baseURL + "/segments",
		"segment_url":                   baseURL + "/segments/{id}",
		"students_url":                  baseURL + "/students",
		"student_url":                   baseURL + "/students/{id}",
		"student_enrollments_url":       baseURL + "/student_enrollments",
		"student_enrollment_url":        baseURL + "/student_enrollments/{id}",
		"teacher_classes_url":           baseURL + "/teacher_classes",
		"teacher_class_url":             baseURL + "/teacher_classes/{id}",
		"users_url":                     baseURL + "/users",
		"user_url":                      baseURL + "/users/{id}",
		"user_accesses_url":             baseURL + "/user_accesses",
		"user_access_url":               baseURL + "/user_accesses/{id}",
		"user_access_profiles_url":      baseURL + "/user_access_profiles",
		"user_access_profile_url":       baseURL + "/user_access_profiles/{id}",
	}

	c.IndentedJSON(http.StatusOK, resources)
}

func APIAuthorizedResources(c *gin.Context) {

}

func APIUnauthorizedResources(c *gin.Context) {

}
