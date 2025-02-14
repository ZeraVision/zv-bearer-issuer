package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"
)

func WriteOut(data any, err error, responseWriter http.ResponseWriter, errorCode ...int) {
	if err != nil {
		failWriteOut(err, responseWriter, errorCode...)
	} else {
		successWriteOut(data, responseWriter)
	}
}

func failWriteOut(err error, responseWriter http.ResponseWriter, errorCode ...int) {
	code := http.StatusBadRequest

	if len(errorCode) > 0 {
		code = errorCode[0]
	} else {
		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "authentication") {
			code = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			code = http.StatusNotFound
		} else if strings.Contains(err.Error(), "invalid api key") {
			code = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "rate limit exceeded") {
			code = http.StatusTooManyRequests
		}
	}

	JsonError(responseWriter, err.Error(), code)
}
func JsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func successWriteOut(data any, responseWriter http.ResponseWriter) {
	doPostProcessing(&data)

	// Identify and log unsupported values
	identifyUnsupportedValues(reflect.ValueOf(data), "")

	jsonData, err := json.Marshal(data)
	if err != nil {
		dataStr := fmt.Sprintf("%+v", data)
		log.Printf("successWriteOut: Error marshaling data to JSON: %v, data: %s", err, dataStr)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	responseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Target")

	responseWriter.Write(jsonData)
}

func identifyUnsupportedValues(value reflect.Value, path string) {
	switch value.Kind() {
	case reflect.Ptr:
		if !value.IsNil() {
			identifyUnsupportedValues(value.Elem(), path)
		}
	case reflect.Interface:
		if !value.IsNil() {
			identifyUnsupportedValues(value.Elem(), path)
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			fieldValue := value.Field(i)
			fieldPath := fmt.Sprintf("%s.%s", path, field.Name)
			identifyUnsupportedValues(fieldValue, fieldPath)
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			mapKeyPath := fmt.Sprintf("%s[%v]", path, key.Interface())
			mapValue := value.MapIndex(key)
			identifyUnsupportedValues(mapValue, mapKeyPath)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			elemValue := value.Index(i)
			identifyUnsupportedValues(elemValue, elemPath)
		}
	case reflect.Float32, reflect.Float64:
		f := value.Float()
		if math.IsInf(f, 0) || math.IsNaN(f) {
			log.Printf("identifyUnsupportedValues: Unsupported value found at %s: %v", path, f)
		}
	}
}

// func successWriteOut(data any, responseWriter http.ResponseWriter) {
// 	doPostProcessing(&data)

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		dataStr := fmt.Sprintf("%+v", data)
// 		log.Printf("successWriteOut: Error marshaling data to JSON: %v, data: %s", err, dataStr)
// 		return
// 	}

// 	responseWriter.Header().Set("Content-Type", "application/json")

// 	// Enable CORS by setting appropriate headers
// 	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
// 	responseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	responseWriter.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Target")

// 	responseWriter.Write(jsonData)
// }

func doPostProcessing(data *any) {
	nullToEmptySlice(data)
	replaceNaNsWithZero(data)
}

func replaceNaNsWithZero(data *any) {
	// Create a copy of the original data
	copy := reflect.New(reflect.TypeOf(*data)).Elem()
	copy.Set(reflect.ValueOf(*data))

	// Process the copy to replace NaNs with zeros
	processValue(reflect.Indirect(copy))

	// Assign the processed copy back to the original data
	*data = copy.Interface()
}

func processValue(v reflect.Value) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			processValue(v.Index(i))
		}
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			processValue(v.Elem())
		}
	case reflect.Float32, reflect.Float64:
		if math.IsNaN(v.Float()) {
			v.SetFloat(0)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// Check if the field can be set, to avoid panic
			if field.CanSet() {
				processValue(field)
			}
		}
	}
}

func nullToEmptySlice(data *any) {
	// Create a copy of the original data
	copy := reflect.New(reflect.TypeOf(*data)).Elem()
	copy.Set(reflect.ValueOf(*data))

	initializeSlices(reflect.Indirect(copy))
	*data = copy.Interface()
}

func initializeSlices(val reflect.Value) {
	//fmt.Printf("Processing Type: %s, Kind: %s\n", val.Type(), val.Kind()) // Diagnostic print

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return
		}

		initializeSlices(val.Elem())

	case reflect.Interface:
		if val.IsNil() {
			return
		}

		initializeSlices(val.Elem())

	case reflect.Struct:
		// Traverse each field of the struct.
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.CanAddr() {
				field = field.Addr()
			}
			initializeSlices(field)
		}

	case reflect.Slice:
		// If nil make empty slice
		if val.IsNil() {
			if val.CanSet() {
				newSlice := reflect.MakeSlice(val.Type(), 0, 0)
				val.Set(newSlice)
			}
		}
	}
}
