package validate

import (
	"album-manager/src/errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	v = validator.New()
)

type validateItem struct {
	Key  string
	Func validator.Func
}

func RegisterValidation() error {
	var err error

	validates := []validateItem{
		{
			Key:  "date_string",
			Func: validateDateFormatString,
		},
		{
			Key:  "alphanumspace",
			Func: validateAlphanumericAndSpace,
		},
	}

	for _, r := range validates {
		err = v.RegisterValidation(r.Key, r.Func)
		if err != nil {
			return err
		}

		log.Printf(`RegisterValidation: "%s" successfully!`, r.Key)
	}

	return nil
}

func getFieldName(req interface{}, namespace string) string {
	t := reflect.TypeOf(req)

	nsSplit := strings.Split(namespace, ".")[1:]
	fieldNames := make([]string, 0)

	var embeddedType reflect.StructField

	for i, n := range nsSplit {
		var embeddedField reflect.StructField
		if i == 0 {
			embeddedField, _ = t.Elem().FieldByName(n)
			embeddedType = embeddedField
		} else {
			embeddedField, _ = embeddedType.Type.FieldByName(n)
		}

		fmt.Println(embeddedField)

		tag := embeddedField.Tag.Get("json")
		if tag == "" {
			tag = embeddedField.Tag.Get("form")
		}

		fieldNames = append(fieldNames, tag)
	}

	return strings.Join(fieldNames, ".")
}

// Valid validates the given struct.
func Valid(dst interface{}) error {
	err := v.Struct(dst)
	if err == nil {
		return nil
	}

	userFacingErrors := make(errors.M)

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := getFieldName(dst, err.Namespace())

		switch err.Tag() {
		case "required":
			userFacingErrors[fieldName] = "This field is required."
		case "min":
			if err.Type().Kind() == reflect.String {
				userFacingErrors[fieldName] =
					fmt.Sprintf("This field must be at least %s characters long.", err.Param())
			} else {
				userFacingErrors[fieldName] =
					fmt.Sprintf("This value does not meet the minimum of %s.", err.Param())
			}
		case "max":
			if err.Type().Kind() == reflect.String {
				userFacingErrors[fieldName] =
					fmt.Sprintf("This field must be less than %s characters long.", err.Param())
			} else {
				userFacingErrors[fieldName] =
					fmt.Sprintf("This value exceeds the maximum of %s.", err.Param())
			}
		case "email":
			userFacingErrors[fieldName] = "This isn't a valid email."
		case "alphanumspace":
			userFacingErrors[fieldName] = "This value isn't Alphanumeric."
		case "alpha":
			userFacingErrors[fieldName] = "This value isn't Alphabet."
		case "date_string":
			userFacingErrors[fieldName] = "Wrong format type (dd/mm/yyyy)."
		default:
			userFacingErrors[fieldName] = "Got some errors"
		}
	}

	return errors.E(errors.Op("payload.Valid"), err, userFacingErrors, http.StatusBadRequest)
}

// ReadValid is equivalent to calling Read followed by Valid.
func ReadValid(dst interface{}, c *gin.Context) error {
	op := errors.Op("utils.validate.ReadValid")

	if err := Read(dst, c); err != nil {
		return errors.E(op, err)
	}

	if err := Valid(dst); err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Read unmarshals the payload from the incoming request to the given struct pointer.
func Read(dst interface{}, c *gin.Context) error {
	if err := c.ShouldBindJSON(dst); err != nil {
		return errors.E(errors.Op("utils.validate.Read"), http.StatusBadRequest, err,
			map[string]string{"message": "Could not decode request body"})
	}

	return nil
}

// Read unmarshals the payload from the incoming request to the given struct pointer.
func ReadQuery(dst interface{}, c *gin.Context) error {
	if err := c.ShouldBindQuery(dst); err != nil {
		return errors.E(errors.Op("utils.validate.ReadQuery"), http.StatusBadRequest, err,
			map[string]string{"message": "Could not decode query params"})
	}

	return nil
}

// ReadValid is equivalent to calling Read followed by Valid.
func ReadQueryValid(dst interface{}, c *gin.Context) error {
	op := errors.Op("utils.validate.ReadValid")

	if err := ReadQuery(dst, c); err != nil {
		return errors.E(op, err)
	}

	if err := Valid(dst); err != nil {
		return errors.E(op, err)
	}

	return nil
}

func lowerFirstLetter(s string) string {
	if r := rune(s[0]); r >= 'A' && r <= 'Z' {
		s = strings.ToLower(string(r)) + s[1:]
	}

	if s[len(s)-2:] == "ID" {
		s = s[:len(s)-2] + "Id"
	}

	return s
}

func validateAlphanumericAndSpace(fl validator.FieldLevel) bool {
	// Regular expression pattern to match alphanumeric characters and spaces
	pattern := "^[a-zA-Z0-9 ]+$"

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	// Retrieve the field value as a string
	value := fl.Field().String()

	// Check if the value matches the regular expression pattern
	return regex.MatchString(value)
}

func validateDateFormatString(fl validator.FieldLevel) bool {
	dob := fl.Field().String()

	if !fl.Field().IsZero() {
		// Regular expression pattern for "dd/mm/yyyy" format
		pattern := `^(0[1-9]|[1-2][0-9]|3[0-1])/(0[1-9]|1[0-2])/\d{4}$`
		match, _ := regexp.MatchString(pattern, dob)

		return match
	}

	return true
}

func validateNullCheck(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Check if the field is an int
	if field.Kind() == reflect.Int {
		// Get the int value
		value := field.Int()

		// If the value is zero, consider it empty
		if value == 0 {
			return true
		}
	}

	// Perform additional validation logic if needed

	return true
}
