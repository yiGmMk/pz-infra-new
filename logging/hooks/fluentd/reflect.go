// Based on `logrus_fluent`
// See https://github.com/evalphobia/logrus_fluent

package fluentd

import (
	"fmt"
	"reflect"
	"strings"
)

// ConvertToValue make map data from struct and tags
func ConvertToValue(p interface{}, tagName string, needReplaceDots bool) interface{} {
	rv := toValue(p)
	switch rv.Kind() {
	case reflect.Struct:
		return converFromStruct(rv.Interface(), tagName, needReplaceDots)
	case reflect.Map:
		return convertFromMap(rv, tagName, needReplaceDots)
	case reflect.Slice:
		return convertFromSlice(rv, tagName, needReplaceDots)
	case reflect.Chan:
		return nil
	case reflect.Invalid:
		return nil
	case reflect.Bool:
		return rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint()
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	case reflect.Complex64, reflect.Complex128:
		return rv.Complex()
	case reflect.String:
		return rv.String()
	default:
		return rv.Interface()
	}
}

func convertFromMap(rv reflect.Value, tagName string, needReplaceDots bool) interface{} {
	result := make(map[string]interface{})
	for _, key := range rv.MapKeys() {
		ks := fmt.Sprint(key.Interface())
		if needReplaceDots {
			ks = replaceDots(ks)
		}
		kv := rv.MapIndex(key)
		result[ks] = ConvertToValue(kv.Interface(), tagName, needReplaceDots)
	}
	return result
}

func convertFromSlice(rv reflect.Value, tagName string, needReplaceDots bool) interface{} {
	var result []interface{}
	for i, max := 0, rv.Len(); i < max; i++ {
		result = append(result, ConvertToValue(rv.Index(i).Interface(), tagName, needReplaceDots))
	}
	return result
}

// converFromStruct converts struct to value
// see: https://github.com/fatih/structs/
func converFromStruct(p interface{}, tagName string, needReplaceDots bool) interface{} {
	t := toType(p)
	// parse Time
	if t.String() == "time.Time" {
		return fmt.Sprint(p)
	}
	result := make(map[string]interface{})
	values := toValue(p)
	for i, max := 0, t.NumField(); i < max; i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue // skip private field
		}
		tag, opts := parseTag(f, tagName)
		if tag == "-" {
			continue // skip `-` tag
		}

		v := values.Field(i)
		if opts.Has("omitempty") && isZero(v) {
			continue // skip zero-value when omitempty option exists in tag
		}
		name := getNameFromTag(f, tagName, needReplaceDots)
		result[name] = ConvertToValue(v.Interface(), tagName, needReplaceDots)
	}
	return result
}

// toValue converts any value to reflect.Value
func toValue(p interface{}) reflect.Value {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// toType converts any value to reflect.Type
func toType(p interface{}) reflect.Type {
	t := reflect.ValueOf(p).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// isZero checks the value is zero-value or not
func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	value := v.Interface()
	return reflect.DeepEqual(value, zero)
}

// getNameFromTag return the value in tag or field name in the struct field
func getNameFromTag(f reflect.StructField, tagName string, needReplaceDots bool) string {
	name, _ := parseTag(f, tagName)
	if name == "" {
		name = f.Name
	}
	if needReplaceDots {
		name = replaceDots(name)
	}
	return name
}

// getTagValues returns tag value of the struct field
func getTagValues(f reflect.StructField, tag string) string {
	return f.Tag.Get(tag)
}

// parseTag returns the first tag value of the struct field
func parseTag(f reflect.StructField, tag string) (string, options) {
	return splitTags(getTagValues(f, tag))
}

// splitTags returns the first tag value and rest slice
func splitTags(tags string) (string, options) {
	res := strings.Split(tags, ",")
	return res[0], res[1:]
}

// TagOptions is wrapper struct for rest tag values
type options []string

// Has checks the value exists in the rest values or not
func (t options) Has(tag string) bool {
	for _, opt := range t {
		if opt == tag {
			return true
		}
	}
	return false
}

func replaceDots(s string) string {
	return strings.Replace(s, ".", "_", -1)
}
