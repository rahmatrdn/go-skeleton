package factory

import (
	"database/sql/driver"
	"reflect"
	"sort"
	"sync"

	"github.com/rahmatrdn/go-skeleton/internal/helper"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm/schema"
)

func GetRows(model any, nonZeroVal bool) []driver.Value {
	keys, m := getKeys(model, nonZeroVal)

	v := make([]driver.Value, 0, len(keys))
	for _, k := range keys {
		v = append(v, m[k])
	}

	return v
}

func GetInsertRows(model any, nonZeroVal bool, customVal map[string]any) []driver.Value {
	keys, m := getKeys(model, nonZeroVal)

	v := make([]driver.Value, 0, len(keys))

loopkey:
	for _, k := range keys {
		for custK := range customVal {
			if k == custK {
				continue loopkey
			}
		}
		v = append(v, m[k])
	}

	for k := range customVal {
		v = append(v, customVal[k])
	}

	return v
}

func GetUpdateRows(model any, nonZeroVal bool, idCol string, idVal any, hasUpdatedAt bool) []driver.Value {
	keys, m := getKeys(model, nonZeroVal)

	v := make([]driver.Value, 0, len(keys))

	for _, k := range keys {
		v = append(v, m[k])
	}

	if nonZeroVal && hasUpdatedAt {
		v = append(v, sqlmock.AnyArg())
	}
	v = append(v, idVal)

	return v
}

func GetUpdateRowsWithUserID(model any, nonZeroVal bool, idVal any, userID int64, hasUpdatedAt bool) []driver.Value {
	v := GetUpdateRows(model, nonZeroVal, "ID", userID, hasUpdatedAt)
	v = append(v, idVal)

	return v
}

func getKeys(model any, nonZeroVal bool) ([]string, map[string]any) {
	m := helper.StructToMap(model, nonZeroVal)

	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys, m
}

func GetCols(model any) []string {
	columns := []string{}
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		panic("failed to create schema")
	}

	keys := make([]string, 0)
	m := map[string]string{}
	for _, field := range s.Fields {
		keys = append(keys, field.Name)
		m[field.Name] = field.DBName
	}
	sort.Strings(keys)

	for _, k := range keys {
		columns = append(columns, m[k])
	}

	return columns
}

func InArray(val interface{}, array interface{}) (found bool) {
	values := reflect.ValueOf(array)

	if reflect.TypeOf(array).Kind() == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}
