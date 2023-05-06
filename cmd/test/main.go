package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Type uint

const (
	_ Type = iota
	Struct
	Array
	Map
	String
	Bool
	Integer
	Float
)

type ElemValue struct {
	Type   Type       `json:"type"`
	Struct *Structure `json:"struct,omitempty"`
	Array  *ElemValue `json:"array,omitempty"`
	Map    *KeyValue  `json:"map,omitempty"`
}

type KeyValue struct {
	Key   Type      `json:"key"`
	Value ElemValue `json:"value"`
}

type Field struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        Type   `json:"type"`
	Optional    bool   `json:"optional,omitempty"`

	Struct *Structure `json:"struct,omitempty"`
	Array  *ElemValue `json:"array,omitempty"`
	Map    *KeyValue  `json:"map,omitempty"`
}

type Structure struct {
	Fields []Field `json:"fields"`
}

func StructureOf[T any]() ([]byte, error) {
	var v T
	structure, err := structureOfStruct(v)
	if err != nil {
		return nil, err
	}

	return json.Marshal(structure)
}

func structureOfStruct(v any) (*Structure, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %q is not a struct", typ)
	}

	var fields []Field
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		f, err := structureOfField(field)
		if err != nil {
			return nil, fmt.Errorf("field: %q, type: %q, err: %w", field.Name, field.Type, err)
		}

		fields = append(fields, f)
	}

	return &Structure{
		Fields: fields,
	}, nil
}

func structureOfField(field reflect.StructField) (Field, error) {
	tags := field.Tag
	name := strings.TrimSpace(tags.Get("name"))
	if name == "" {
		return Field{}, fmt.Errorf("no name provided")
	}

	f := Field{
		Name:        name,
		Description: strings.TrimSpace(tags.Get("desc")),
	}

	fieldType := field.Type
	fieldKind := fieldType.Kind()
	if fieldKind == reflect.Pointer {
		f.Optional = true
		fieldType = fieldType.Elem()
		fieldKind = fieldType.Kind()
	}

	switch fieldKind {
	case reflect.String:
		f.Type = String
	case reflect.Bool:
		f.Type = Bool
	case reflect.Int64:
		f.Type = Integer
	case reflect.Float64:
		f.Type = Float
	case reflect.Struct:
		f.Type = Struct

		fv := reflect.Zero(fieldType)
		if !fv.CanInterface() {
			return Field{}, errors.New("can't instantiate type")
		}

		fs, err := structureOfStruct(fv.Interface())
		if err != nil {
			return Field{}, fmt.Errorf("structure of, err: %w", err)
		}

		f.Struct = fs
	case reflect.Slice:
		f.Type = Array

		elemType := fieldType.Elem()
		elemValue, err := structureOfElem(elemType)
		if err != nil {
			return Field{}, fmt.Errorf("array element type: %q, err: %w", elemType, err)
		}

		f.Array = elemValue
	case reflect.Map:
		f.Type = Map

		f.Map = &KeyValue{}

		var err error
		f.Map.Key, err = structureOfKey(fieldType.Key())
		if err != nil {
			return Field{}, err
		}

		elemType := fieldType.Elem()
		elemValue, err := structureOfElem(elemType)
		if err != nil {
			return Field{}, fmt.Errorf("map value type: %q, err: %w", elemType, err)
		}

		f.Map.Value = *elemValue
	default:
		return Field{}, fmt.Errorf("unsoppurted field kind %q", fieldKind)
	}

	return f, nil
}

func structureOfKey(typ reflect.Type) (Type, error) {
	var keyType Type

	keyKind := typ.Kind()
	switch keyKind {
	case reflect.String:
		keyType = String
	case reflect.Bool:
		keyType = Bool
	case reflect.Int64:
		keyType = Integer
	case reflect.Float64:
		keyType = Float
	default:
		return 0, fmt.Errorf("unsoppurted map key kind %q", keyKind)
	}

	return keyType, nil
}

func structureOfElem(typ reflect.Type) (*ElemValue, error) {
	v := &ElemValue{}

	elemKind := typ.Kind()
	switch elemKind {
	case reflect.String:
		v.Type = String
	case reflect.Bool:
		v.Type = Bool
	case reflect.Int64:
		v.Type = Integer
	case reflect.Float64:
		v.Type = Float
	case reflect.Struct:
		v.Type = Struct

		fv := reflect.Zero(typ)
		if !fv.CanInterface() {
			return nil, errors.New("can't instantiate type")
		}

		fs, err := structureOfStruct(fv.Interface())
		if err != nil {
			return nil, fmt.Errorf("structure of, err: %w", err)
		}

		v.Struct = fs
	case reflect.Slice:
		v.Type = Array

		elemType := typ.Elem()
		vv, err := structureOfElem(elemType)
		if err != nil {
			return nil, fmt.Errorf("array type: %q, err: %w", elemType, err)
		}

		v.Array = vv
	case reflect.Map:
		v.Type = Map

		v.Map = &KeyValue{}

		var err error
		v.Map.Key, err = structureOfKey(typ.Key())
		if err != nil {
			return nil, err
		}

		elemType := typ.Elem()
		vv, err := structureOfElem(elemType)
		if err != nil {
			return nil, fmt.Errorf("map value type: %q, err: %w", elemType, err)
		}

		v.Map.Value = *vv
	default:
		return nil, fmt.Errorf("unsoppurted array element kind %q", elemKind)
	}

	return v, nil
}

type Button struct {
	Name string `json:"name" name:"Name"`
}

type Keyboard struct {
	ButtonWidth float64   `json:"button_width" name:"Button width"`
	Buttons     []Button  `json:"buttons"      name:"Buttons"`
	Indexes     []int64   `json:"indexes"      name:"Indexes"`
	Points      [][]int64 `json:"points"       name:"Points"`
}

type SendMsg struct {
	ChatID       int64                       `json:"chat_id"      name:"Chat ID"`
	Text         string                      `json:"text"         name:"Text" desc:"Message text to send"`
	Keyboard     *Keyboard                   `json:"keyboard"     name:"Keyboard"`
	Languages    map[string]bool             `json:"languages"    name:"Languages" desc:"Supported languages"`
	Users        map[int64]struct{}          `json:"users"        name:"Users"`
	UserSettings map[int64]map[string]string `json:"userSettings" name:"User settings"`
	SliceOfMaps  []map[bool]string           `json:"sliceOfMaps"  name:"som"`
}

func main() {
	start := time.Now()

	data, err := StructureOf[SendMsg]()
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Since(start))
	fmt.Println(string(data))
}
