package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

const (
	tag            = "validate"
	logicSeparator = "|"
	valueSeparator = ","
)

var (
	ErrStruct = errors.New("not a struct")

	ErrStrLength       = errors.New("string 'length' validation problem")
	ErrStrRegexp       = errors.New("string 'regexp' validation problem")
	ErrStrIn           = errors.New("string 'in' validation problem")
	ErrStrLengthFormat = errors.New("string 'length' format problem")
	ErrStrRegexpFormat = errors.New("string 'regexp' format problem")
	ErrStrInFormat     = errors.New("string 'in' format problem")

	ErrIntMax       = errors.New("int 'max' validation problem")
	ErrIntMin       = errors.New("int 'min' validation problem")
	ErrIntIn        = errors.New("int 'in' validation problem")
	ErrIntMaxFormat = errors.New("int 'max' format problem")
	ErrIntMinFormat = errors.New("int 'min' format problem")
	ErrIntInFormat  = errors.New("int 'in' format problem")

	ErrType = errors.New("field type validation problem")
)

func (v ValidationErrors) Error() string {
	var result strings.Builder
	for i, err := range v {
		result.WriteString(err.Field)
		result.WriteString(" field: ")
		result.WriteString(err.Err.Error())
		if i != len(v)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}

func Validate(vi interface{}) error {
	var validationErrors ValidationErrors
	v := reflect.ValueOf(vi)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validation error: %w: expected a struct, got %T", ErrStruct, vi)
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if tagValue, ok := fieldType.Tag.Lookup(tag); ok && len(tagValue) > 0 && fieldValue.CanInterface() {
			rules, err := getRules(tagValue, fieldValue.Kind())
			if err != nil {
				return err
			}
			for _, rule := range rules {
				validationErrors = append(validationErrors, rule.Validate(fieldType.Name, fieldValue)...)
			}
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func getRules(tagValue string, fieldKind reflect.Kind) ([]ValidationRule, error) {
	rules := strings.Split(tagValue, logicSeparator)
	validationRules := make([]ValidationRule, 0, len(rules))
	for _, ruleStr := range rules {
		var err error
		var rule ValidationRule
		switch {
		case strings.HasPrefix(ruleStr, "len:"):
			rule, err = NewStringLenRule(ruleStr)
			validationRules = append(validationRules, rule)
		case strings.HasPrefix(ruleStr, "regexp:"):
			rule, err = NewStringRegexpRule(ruleStr)
			validationRules = append(validationRules, rule)
		case strings.HasPrefix(ruleStr, "min:"):
			rule, err = NewIntMinRule(ruleStr)
			validationRules = append(validationRules, rule)
		case strings.HasPrefix(ruleStr, "max:"):
			rule, err = NewIntMaxRule(ruleStr)
			validationRules = append(validationRules, rule)
		case strings.HasPrefix(ruleStr, "in:"):
			if fieldKind == reflect.Int {
				rule, err = NewIntInRule(ruleStr)
				validationRules = append(validationRules, rule)
			}
			if fieldKind == reflect.String {
				rule, err = NewStringInRule(ruleStr)
				validationRules = append(validationRules, rule)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return validationRules, nil
}

type validator = func(fieldName string, fieldValue reflect.Value) *ValidationError

type TypeValidationRule interface {
	Validate(fieldName string, fieldValue reflect.Value, validationLogic validator) ValidationErrors
}

type typeRule struct {
	typeKind reflect.Kind
}

func (rule typeRule) Validate(fieldName string, fieldValue reflect.Value, validate validator) ValidationErrors {
	switch fieldValue.Kind() { //nolint:exhaustive
	case rule.typeKind:
		validationError := validate(fieldName, fieldValue)
		if validationError != nil {
			return []ValidationError{*validationError}
		}
		return nil
	case reflect.Slice, reflect.Array:
		errs := ValidationErrors{}
		if fieldValue.IsNil() {
			return errs
		}
		for i := 0; i < fieldValue.Len(); i++ {
			elem := fieldValue.Index(i)
			if elem.Kind() != rule.typeKind {
				errs = append(errs, ValidationError{
					fieldName,
					fmt.Errorf("%w: field value type is not the %s type", ErrType, rule.typeKind.String()),
				})
			}
			if err := validate(fieldName, elem); err != nil {
				errs = append(errs, *err)
			}
		}
		return errs
	default:
		return []ValidationError{{
			fieldName,
			fmt.Errorf("%w: field is not %s or %s type", ErrType, rule.typeKind.String(), rule.typeKind.String()),
		}}
	}
}

func getRuleValue(ruleStr string, ruleStrPrefix string) string {
	if len(ruleStr) == 0 {
		return ""
	}
	if strings.HasPrefix(ruleStr, ruleStrPrefix) {
		ruleStrParts := strings.SplitN(ruleStr, ":", 2)
		if len(ruleStrParts) > 1 {
			return ruleStrParts[1]
		}
	}
	return ""
}

type ValidationRule interface {
	Validate(fieldName string, fieldValue reflect.Value) ValidationErrors
}

type stringLenRule struct {
	typeRule TypeValidationRule
	len      int
}

func (rule stringLenRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if len(value.String()) != rule.len {
			return &ValidationError{name, fmt.Errorf("%w: field length not equals to %d", ErrStrLength, rule.len)}
		}
		return nil
	})
}

func NewStringLenRule(ruleStr string) (ValidationRule, error) {
	ruleValue := getRuleValue(ruleStr, "len:")
	length, err := strconv.Atoi(ruleValue)
	if err != nil {
		return nil, ErrStrLengthFormat
	}
	return stringLenRule{typeRule{reflect.String}, length}, nil
}

type stringRegexpRule struct {
	typeRule TypeValidationRule
	regexp   *regexp.Regexp
}

func (rule stringRegexpRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if !rule.regexp.MatchString(value.String()) {
			return &ValidationError{
				name,
				fmt.Errorf("%w: field is not matching regular expression", ErrStrRegexp),
			}
		}
		return nil
	})
}

func NewStringRegexpRule(ruleStr string) (ValidationRule, error) {
	ruleValue := getRuleValue(ruleStr, "regexp:")
	r, err := regexp.Compile(ruleValue)
	if err != nil {
		return nil, ErrStrRegexpFormat
	}
	return stringRegexpRule{typeRule{reflect.String}, r}, nil
}

type stringInRule struct {
	typeRule        TypeValidationRule
	availableValues []string
}

func (rule stringInRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		for _, availableStr := range rule.availableValues {
			if availableStr == value.String() {
				return nil
			}
		}
		return &ValidationError{
			name,
			fmt.Errorf("%w field value is not matching: %v", ErrStrIn, rule.availableValues),
		}
	})
}

func NewStringInRule(ruleStr string) (ValidationRule, error) {
	ruleValue := getRuleValue(ruleStr, "in:")
	valuesList := strings.Split(ruleValue, valueSeparator)
	if len(valuesList) == 0 {
		return nil, ErrStrInFormat
	}
	return stringInRule{typeRule{reflect.String}, valuesList}, nil
}

type intMinRule struct {
	typeRule TypeValidationRule
	min      int64
}

func (rule intMinRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if value.Int() < rule.min {
			return &ValidationError{
				name,
				fmt.Errorf("%w: value is lower than min", ErrIntMin),
			}
		}
		return nil
	})
}

func NewIntMinRule(ruleStr string) (ValidationRule, error) {
	ruleValue := getRuleValue(ruleStr, "min:")
	min, err := strconv.ParseInt(ruleValue, 10, 0)
	if err != nil {
		return nil, ErrIntMinFormat
	}
	return intMinRule{typeRule{reflect.Int}, min}, nil
}

type intMaxRule struct {
	typeRule TypeValidationRule
	max      int64
}

func (rule intMaxRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if value.Int() > rule.max {
			return &ValidationError{
				name,
				fmt.Errorf("%w: value is bigger than max", ErrIntMax),
			}
		}
		return nil
	})
}

func NewIntMaxRule(ruleStr string) (ValidationRule, error) {
	ruleValue := getRuleValue(ruleStr, "max:")
	max, err := strconv.ParseInt(ruleValue, 10, 0)
	if err != nil {
		return nil, ErrIntMaxFormat
	}
	return intMaxRule{typeRule{reflect.Int}, max}, nil
}

type intInRule struct {
	typeRule            TypeValidationRule
	availableValuesList []int64
}

func (rule intInRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		for _, availableValue := range rule.availableValuesList {
			if availableValue == value.Int() {
				return nil
			}
		}
		return &ValidationError{
			name,
			fmt.Errorf("%w: value is not matching: %v", ErrIntIn, rule.availableValuesList),
		}
	})
}

func NewIntInRule(ruleStr string) (ValidationRule, error) {
	ruleValue := getRuleValue(ruleStr, "in:")
	strValues := strings.Split(ruleValue, valueSeparator)
	valuesList := make([]int64, 0, len(strValues))
	for _, strValue := range strValues {
		parsedInt, err := strconv.ParseInt(strValue, 10, 0)
		if err != nil {
			return nil, ErrIntInFormat
		}
		valuesList = append(valuesList, parsedInt)
	}
	if len(valuesList) == 0 {
		return nil, ErrIntInFormat
	}
	return intInRule{typeRule{reflect.Int}, valuesList}, nil
}
