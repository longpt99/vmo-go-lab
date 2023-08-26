package object

import (
	"album-manager/src/errors"
	"net/http"
	"reflect"
)

func MergeStructIntoModel(dest, body interface{}) error {
	bodyReq := reflect.ValueOf(body)
	paramsValue := reflect.ValueOf(dest)

	if bodyReq.Kind() != reflect.Ptr || paramsValue.Kind() != reflect.Ptr {
		return errors.E(errors.Op("MergeStructIntoModel"), http.StatusBadRequest, "params must be pointer")
	}

	bodyReq = bodyReq.Elem()
	paramsValue = paramsValue.Elem()
	// Get the reflect.Value of the createReq variable

	// Iterate over the fields of createReq
	for i := 0; i < bodyReq.NumField(); i++ {
		// Get the field name and value
		fieldName := bodyReq.Type().Field(i).Name
		fieldValue := bodyReq.Field(i)

		// Set the corresponding field in params using type assertion
		paramsValue.FieldByName(fieldName).Set(fieldValue)
	}

	return nil
}
