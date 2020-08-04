package main

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

var typeMap = map[string]interface{}{}

// Party ...
type Party struct {
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Attendees []string  `json:"attendees" validate:"required,gt=0"`
}

// MovieParty ...
type MovieParty struct {
	Party
	Movie   string `json:"movie" validate:"required"`
	Rating  string `json:"rating" validate:"required,oneof=G PG PG-13 R NC-17"`
	Runtime int    `json:"runtime" validate:"required,gt=30"`
}

// PoolParty ...
type PoolParty struct {
	Party
	WaterTemp int `json:"water_temp" validate:"required"`
}

// DinnerParty ...
type DinnerParty struct {
	Party
	Dinner  string `json:"dinner" validate:"required"`
	Dessert string `json:"dessert" validate:""`
}

// FieldSpec is a field description
type FieldSpec struct {
	Typ      string   `json:"type"`
	Checks   []string `json:"checks"`
	Required bool     `json:"required"`
	List     bool     `json:"list"`
}

func getTypeData(t reflect.Type) map[string]FieldSpec {
	m := make(map[string]FieldSpec)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			am := getTypeData(f.Type)
			for k, v := range am {
				m[k] = v
			}
		} else {
			typeName, array := getTypeString(f.Type.String())
			checks, required := getChecks(f)
			if len(checks) == 0 {
				checks = nil
			}
			//translate the checks to something more meaningful for lists or
			//those that reference other fields
			for i, check := range checks {
				if check == "" {
					continue
				}
				checkParts := strings.Split(check, "=")
				checkName := checkParts[0]
				checkVal := checkParts[1]
				if strings.HasSuffix(checkName, "field") {
					realField, _ := t.FieldByName(checkVal)
					checkVal = getJSONName(realField)
				}
				checks[i] = checkName + "=" + checkVal
			}
			//todo massage gt for list to be lengt
			//todo massage *field to be json name instead of field name (womp womp)
			m[getJSONName(f)] = FieldSpec{
				Typ:      typeName,
				List:     array,
				Checks:   checks,
				Required: required,
			}
		}
	}
	return m
}

func getJSONName(f reflect.StructField) string {
	jsonParts := strings.Split(f.Tag.Get("json"), ",")
	return jsonParts[0]
}

func getChecks(f reflect.StructField) ([]string, bool) {
	req := false
	checkParts := strings.Split(f.Tag.Get("validate"), ",")
	for i := 0; i < len(checkParts); i++ {
		if checkParts[i] == "required" {
			req = true
			checkParts = append(checkParts[:i], checkParts[i+1:]...)
		}
	}
	return checkParts, req
}

func getTypeString(ts string) (string, bool) {
	if ts == "time.Time" {
		return "RFC 3339", false
	}
	if strings.HasPrefix(ts, "[]") {
		return ts[2:], true
	}
	return ts, false
}

func init() {
	for _, t := range []reflect.Type{
		reflect.TypeOf((*MovieParty)(nil)).Elem(),
		reflect.TypeOf((*PoolParty)(nil)).Elem(),
		reflect.TypeOf((*DinnerParty)(nil)).Elem(),
	} {
		typeMap[t.Name()] = getTypeData(t)
	}
}
