package main

import (
	"encoding/json"
	"fmt"

	"github.com/zephinzer/go-config"
)

type MyConfiguration struct {
	OptionalBool                               *bool
	OptionalBoolWithCustomEnv                  *bool `env:"CUSTOM_BOOL"`
	OptionalInt                                *int
	OptionalString                             *string
	OptionalStringSlice                        *[]string
	RequiredBool                               bool
	RequiredBoolWithDefault                    bool `default:"true"`
	RequiredInt                                int
	RequiredIntWithDefault                     int `default:"-12345"`
	RequiredString                             string
	RequiredStringWithDefault                  string `default:"required-string"`
	RequiredStringSlice                        []string
	RequiredStringSliceWithDefault             []string `default:"required,string,slice"`
	RequiredStringSliceWithDefaultAndDelimiter []string `default:"required|string|slice" delimiter:"|"`
}

func main() {
	var configuration MyConfiguration
	if err := config.Load(&configuration); err != nil {
		panic(err)
	}
	jsonifiedConfig, err := json.MarshalIndent(configuration, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonifiedConfig))
}
