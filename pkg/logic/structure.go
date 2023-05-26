package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func StructureOf[T any]() ([]byte, error) {
	var v T
	structure, err := structureOfStruct(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(structure)
}

func MustStructureOf[T any]() []byte {
	structure, err := StructureOf[T]()
	if err != nil {
		var v T
		panic(fmt.Errorf("structure of %T: %w", v, err))
	}
	return structure
}

type valueType uint

const (
	_ valueType = iota
	valueTypeStruct
	valueTypeArray
	valueTypeMap
	valueTypeString
	valueTypeBool
	valueTypeInteger
	valueTypeFloat
)

type elemValue struct {
	Type   valueType       `json:"type"`
	Struct *structureValue `json:"struct,omitempty"`
	Array  *elemValue      `json:"array,omitempty"`
	Map    *keyValue       `json:"map,omitempty"`
}

type keyValue struct {
	Key   valueType `json:"key"`
	Value elemValue `json:"value"`
}

type fieldValue struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Type        valueType `json:"type"`
	Optional    bool      `json:"optional,omitempty"`

	Struct *structureValue `json:"struct,omitempty"`
	Array  *elemValue      `json:"array,omitempty"`
	Map    *keyValue       `json:"map,omitempty"`
}

type structureValue struct {
	Fields []fieldValue `json:"fields"`
}

func structureOfStruct(v any) (*structureValue, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %q is not a struct", typ)
	}

	var fields []fieldValue
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

	return &structureValue{
		Fields: fields,
	}, nil
}

func structureOfField(field reflect.StructField) (fieldValue, error) {
	tags := field.Tag
	name := strings.TrimSuffix(tags.Get("json"), ",omitempty")
	if name == "" {
		return fieldValue{}, fmt.Errorf("no name provided")
	}

	f := fieldValue{
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
		f.Type = valueTypeString
	case reflect.Bool:
		f.Type = valueTypeBool
	case reflect.Int64:
		f.Type = valueTypeInteger
	case reflect.Float64:
		f.Type = valueTypeFloat
	case reflect.Struct:
		f.Type = valueTypeStruct

		fv := reflect.Zero(fieldType)
		if !fv.CanInterface() {
			return fieldValue{}, errors.New("can't instantiate type")
		}

		fs, err := structureOfStruct(fv.Interface())
		if err != nil {
			return fieldValue{}, fmt.Errorf("structure of, err: %w", err)
		}

		f.Struct = fs
	case reflect.Slice:
		f.Type = valueTypeArray

		elemType := fieldType.Elem()
		value, err := structureOfElem(elemType)
		if err != nil {
			return fieldValue{}, fmt.Errorf("array element type: %q, err: %w", elemType, err)
		}

		f.Array = value
	case reflect.Map:
		f.Type = valueTypeMap

		f.Map = &keyValue{}

		var err error
		f.Map.Key, err = structureOfKey(fieldType.Key())
		if err != nil {
			return fieldValue{}, err
		}

		elemType := fieldType.Elem()
		value, err := structureOfElem(elemType)
		if err != nil {
			return fieldValue{}, fmt.Errorf("map value type: %q, err: %w", elemType, err)
		}

		f.Map.Value = *value
	default:
		return fieldValue{}, fmt.Errorf("unsupported field kind %q", fieldKind)
	}

	return f, nil
}

func structureOfKey(typ reflect.Type) (valueType, error) {
	var keyType valueType

	keyKind := typ.Kind()
	switch keyKind {
	case reflect.String:
		keyType = valueTypeString
	case reflect.Bool:
		keyType = valueTypeBool
	case reflect.Int64:
		keyType = valueTypeInteger
	case reflect.Float64:
		keyType = valueTypeFloat
	default:
		return 0, fmt.Errorf("unsupported map key kind %q", keyKind)
	}

	return keyType, nil
}

func structureOfElem(typ reflect.Type) (*elemValue, error) {
	v := &elemValue{}

	elemKind := typ.Kind()
	switch elemKind {
	case reflect.String:
		v.Type = valueTypeString
	case reflect.Bool:
		v.Type = valueTypeBool
	case reflect.Int64:
		v.Type = valueTypeInteger
	case reflect.Float64:
		v.Type = valueTypeFloat
	case reflect.Struct:
		v.Type = valueTypeStruct

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
		v.Type = valueTypeArray

		elemType := typ.Elem()
		value, err := structureOfElem(elemType)
		if err != nil {
			return nil, fmt.Errorf("array type: %q, err: %w", elemType, err)
		}

		v.Array = value
	case reflect.Map:
		v.Type = valueTypeMap

		v.Map = &keyValue{}

		var err error
		v.Map.Key, err = structureOfKey(typ.Key())
		if err != nil {
			return nil, err
		}

		elemType := typ.Elem()
		value, err := structureOfElem(elemType)
		if err != nil {
			return nil, fmt.Errorf("map value type: %q, err: %w", elemType, err)
		}

		v.Map.Value = *value
	default:
		return nil, fmt.Errorf("unsupported array element kind %q", elemKind)
	}

	return v, nil
}
