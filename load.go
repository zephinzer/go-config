package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultStringSliceDelimiter = ","
	ErrorLoadPrereqs            = 1 << iota
	ErrorLoadNotFound           = 1 << iota
	ErrorLoadInvalidType        = 1 << iota
	ErrorLoadInvalidValue       = 1 << iota
)

type LoadErrors []LoadError

func (e LoadErrors) GetCode() int {
	code := 0
	for _, err := range e {
		code |= err.Code
	}
	return code
}

func (e LoadErrors) GetMessage() string {
	messages := []string{}
	for _, err := range e {
		messages = append(messages, err.Message)
	}
	return fmt.Sprintf("['%s']", strings.Join(messages, "', '"))
}

func (e LoadErrors) Error() string {
	codes := 0
	messages := []string{}
	for _, err := range e {
		codes |= err.Code
		messages = append(messages, err.Message)
	}
	return fmt.Sprintf("Load/err[%v]: ['%s']", codes, strings.Join(messages, "', '"))
}

type LoadError struct {
	Code    int
	Message string
}

func (e LoadError) Error() string {
	return fmt.Sprintf("Load/err[%v]: %s", e.Code, e.Message)
}

func Load(config interface{}) error {
	errors := LoadErrors{}

	c := newConfiguration(config)
	if !c.IsPointer() {
		errors = append(errors, LoadError{ErrorLoadPrereqs, "failed to receive a valid pointer"})
	}
	if !c.IsStruct() {
		errors = append(errors, LoadError{ErrorLoadPrereqs, "failed to receive a valid struct"})
	}

	if len(errors) > 0 {
		return errors
	}

	for _, field := range c.Fields {
		environmentKey := field.GetEnvironmentKey()
		environmentValue, isEnvironmentDefined := os.LookupEnv(environmentKey)
		defaultValue := field.GetDefaultValue()
		fieldType := field.Type.String()
		switch fieldType {
		case "[]string":
			var stringSliceValue []string
			delimiter, found := field.Tag.Lookup(StructTagKeyDelimiter)
			if !found {
				delimiter = DefaultStringSliceDelimiter
			}
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadError{
					ErrorLoadNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (string)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			stringValue = strings.Trim(stringValue, delimiter)
			stringSliceValue = strings.Split(stringValue, delimiter)
			field.SetStringSlice(stringSliceValue)
		case "*[]string":
			var stringSliceValue []string
			delimiter, found := field.Tag.Lookup("delimiter")
			if !found {
				delimiter = DefaultStringSliceDelimiter
			}
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			stringValue = strings.Trim(stringValue, delimiter)
			stringSliceValue = strings.Split(stringValue, delimiter)
			field.SetStringSlicePointer(stringSliceValue)
		case "string":
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadError{
					ErrorLoadNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (string)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			field.SetString(stringValue)
		case "*string":
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			field.SetStringPointer(stringValue)
		case "bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				if boolValue, err = strconv.ParseBool(*defaultValue); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadError{
					ErrorLoadNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (bool)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				if boolValue, err = strconv.ParseBool(environmentValue); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetBool(boolValue)
		case "*bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				if boolValue, err = strconv.ParseBool(*defaultValue); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				if boolValue, err = strconv.ParseBool(environmentValue); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetBoolPointer(boolValue)
		case "int":
			var intValue int64
			var err error
			if defaultValue != nil {
				if intValue, err = strconv.ParseInt(*defaultValue, 10, 0); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadError{
					ErrorLoadNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (int)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				if intValue, err = strconv.ParseInt(environmentValue, 10, 0); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetInt(int(intValue))
		case "*int":
			var intValue int64
			var err error
			if defaultValue != nil {
				if intValue, err = strconv.ParseInt(*defaultValue, 10, 0); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				if intValue, err = strconv.ParseInt(environmentValue, 10, 0); err != nil {
					errors = append(errors, LoadError{
						ErrorLoadInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetIntPointer(int(intValue))
		default:
			errors = append(errors, LoadError{
				ErrorLoadInvalidType,
				fmt.Sprintf("failed to load '%s' (via \"${%s}\") of type '%s'", field.Name, environmentKey, field.Type.String()),
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
