package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoadTest struct {
	suite.Suite
}

func TestLoad(t *testing.T) {
	suite.Run(t, &LoadTest{})
}

func (s LoadTest) BeforeTest(string, string) {
	os.Setenv("TEST_ENV_INVALID", "nope")
	os.Setenv("TEST_ENVKEY", "1")
	os.Setenv("TEST_BASE", "1")
}

func (s LoadTest) AfterTest(string, string) {
	os.Unsetenv("TEST_BASE")
	os.Unsetenv("TEST_ENVKEY")
	os.Unsetenv("TEST_ENV_INVALID")
}

func (s LoadTest) TestLoadErrors() {
	errs := LoadErrors{
		{1, "expected message 1"},
		{2, "expected message 2"},
		{4, "expected message 3"},
	}
	s.Equal(1|2|4, errs.GetCode())
	messages := []string{errs.GetMessage(), errs.Error()}
	for _, message := range messages {
		s.Contains(message, "expected message 1")
		s.Contains(message, "expected message 2")
		s.Contains(message, "expected message 3")
	}
}

func (s LoadTest) TestLoadError() {
	err := LoadError{1, "expected message"}
	message := err.Error()
	s.Contains(message, "1")
	s.Contains(message, "expected message")
}

func (s LoadTest) TestLoad_validation() {
	type testStruct struct{}
	err := Load(testStruct{})
	s.NotNil(err)
	s.Equal(ErrorLoadPrereqs, err.(LoadErrors).GetCode())
	s.Contains(err.Error(), "valid pointer")

	var testString string
	err = Load(&testString)
	s.NotNil(err)
	s.Equal(ErrorLoadPrereqs, err.(LoadErrors).GetCode())
	s.Contains(err.Error(), "valid struct")
	err = Load(testString)
	s.NotNil(err)
	s.Equal(ErrorLoadPrereqs, err.(LoadErrors).GetCode())
	s.Contains(err.Error(), "valid struct")

	type testInvalidTypeStruct struct {
		Float float32
	}
	err = Load(&testInvalidTypeStruct{})
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidType, err.(LoadErrors).GetCode())
	s.Contains(err.Error(), "failed to load")
	s.Contains(err.Error(), "type 'float32'")
}

func (s LoadTest) TestLoad_multipleErrors() {
	type testStruct struct {
		Bool   bool
		Int    int `default:"not an int"`
		String string
		Float  float32
	}
	err := Load(&testStruct{})
	s.NotNil(err)
	detailedError, ok := err.(LoadErrors)
	s.True(ok)
	code := detailedError.GetCode()
	s.Equal(ErrorLoadNotFound, code&ErrorLoadNotFound)
	s.Equal(ErrorLoadInvalidValue, code&ErrorLoadInvalidValue)
	s.Equal(ErrorLoadInvalidType, code&ErrorLoadInvalidType)
	messages := []string{detailedError.GetMessage(), detailedError.Error()}
	for _, message := range messages {
		s.Contains(message, `via "${BOOL}" (bool)`)
		s.Contains(message, `parse 'not an int' as an int`)
		s.Contains(message, `via "${STRING}" (string)`)
		s.Contains(message, "of type 'float32'")
	}
}

func (s LoadTest) TestLoad_Bool() {
	type testStruct struct {
		TestBase        bool
		Optional        *bool
		OptionalDefault *bool `default:"true"`
		Default         bool  `default:"true"`
		Env             bool  `env:"TEST_ENVKEY"`
	}
	instance := testStruct{}
	s.NotPanics(func() { Load(&instance) })
	s.Equal(true, instance.TestBase)
	s.Nil(instance.Optional)
	s.Equal(true, *instance.OptionalDefault)
	s.Equal(true, instance.Default)
	s.Equal(true, instance.Env)
}

func (s LoadTest) TestLoad_Bool_notFoundError() {
	type testStruct struct {
		Error bool
	}
	instance := testStruct{}
	err := Load(&instance)
	s.NotNil(err)
	s.Equal(ErrorLoadNotFound, err.(LoadErrors).GetCode())
}

func (s LoadTest) TestLoad_Bool_parseError() {
	type testStruct struct {
		Error bool `default:"nope"`
	}
	instance := testStruct{}
	err := Load(&instance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

	type testStructPointer struct {
		Error *bool `default:"nope"`
	}
	pointerInstance := testStructPointer{}
	err = Load(&pointerInstance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

	type testStructInvalidEnv struct {
		Error bool `env:"TEST_ENV_INVALID"`
	}
	invalidEnvInstance := testStructInvalidEnv{}
	err = Load(&invalidEnvInstance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

	type testStructPointerInvalidEnv struct {
		Error *bool `env:"TEST_ENV_INVALID"`
	}
	invalidEnvPointerInstance := testStructPointerInvalidEnv{}
	err = Load(&invalidEnvPointerInstance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

}

func (s LoadTest) TestLoad_Int() {
	type testStruct struct {
		TestBase        int
		Optional        *int
		OptionalDefault *int `default:"2"`
		Default         int  `default:"3"`
		Env             int  `env:"TEST_ENVKEY"`
	}
	instance := testStruct{}
	s.NotPanics(func() { Load(&instance) })
	s.Equal(1, instance.TestBase)
	s.Nil(instance.Optional)
	s.Equal(2, *instance.OptionalDefault)
	s.Equal(3, instance.Default)
	s.Equal(1, instance.Env)
}

func (s LoadTest) TestLoad_Int_notFoundError() {
	type testStruct struct {
		Error int
	}
	instance := testStruct{}
	err := Load(&instance)
	s.NotNil(err)
	s.Equal(ErrorLoadNotFound, err.(LoadErrors).GetCode())
}

func (s LoadTest) TestLoad_Int_parseError() {
	type testStruct struct {
		Error int `default:"nope"`
	}
	instance := testStruct{}
	err := Load(&instance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

	type testStructPointer struct {
		Error *int `default:"nope"`
	}
	pointerInstance := testStructPointer{}
	err = Load(&pointerInstance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

	type testStructInvalidEnv struct {
		Error int `env:"TEST_ENV_INVALID"`
	}
	invalidEnvInstance := testStructInvalidEnv{}
	err = Load(&invalidEnvInstance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())

	type testStructPointerInvalidEnv struct {
		Error *int `env:"TEST_ENV_INVALID"`
	}
	invalidEnvPointerInstance := testStructPointerInvalidEnv{}
	err = Load(&invalidEnvPointerInstance)
	s.NotNil(err)
	s.Equal(ErrorLoadInvalidValue, err.(LoadErrors).GetCode())
}

func (s LoadTest) TestLoad_String() {
	type testStruct struct {
		TestBase        string
		Optional        *string
		OptionalDefault *string `default:"hi"`
		Default         string  `default:"hello"`
		Env             string  `env:"TEST_ENVKEY"`
		EnvPointer      *string `env:"TEST_ENVKEY"`
	}
	instance := testStruct{}
	s.NotPanics(func() { Load(&instance) })
	s.Equal("1", instance.TestBase)
	s.Nil(instance.Optional)
	s.Equal("hi", *instance.OptionalDefault)
	s.Equal("hello", instance.Default)
	s.Equal("1", instance.Env)
	s.Equal("1", *instance.EnvPointer)
}

func (s LoadTest) TestLoad_String_notFoundError() {
	type testStruct struct {
		Error string
	}
	instance := testStruct{}
	err := Load(&instance)
	s.NotNil(err)
	s.Equal(ErrorLoadNotFound, err.(LoadErrors).GetCode())
}

func (s LoadTest) TestLoad_StringSlice() {
	type testStruct struct {
		TestBase                 []string
		Optional                 *[]string
		OptionalDefault          *[]string `default:"hola,mundo"`
		OptionalDefaultDelimiter *[]string `default:"hi hi" delimiter:" "`
		Default                  []string  `default:"hello"`
		Env                      []string  `env:"TEST_ENVKEY"`
		EnvPointer               *[]string `env:"TEST_ENVKEY"`
	}
	instance := testStruct{}
	s.NotPanics(func() { Load(&instance) })
	s.EqualValues([]string{"1"}, instance.TestBase)
	s.Nil(instance.Optional)
	s.EqualValues([]string{"hola", "mundo"}, *instance.OptionalDefault)
	s.EqualValues([]string{"hi", "hi"}, *instance.OptionalDefaultDelimiter)
	s.EqualValues([]string{"hello"}, instance.Default)
	s.EqualValues([]string{"1"}, instance.Env)
	s.EqualValues([]string{"1"}, *instance.EnvPointer)
}

func (s LoadTest) TestLoad_StringSlice_notFoundError() {
	type testStruct struct {
		Error []string
	}
	instance := testStruct{}
	err := Load(&instance)
	s.NotNil(err)
	s.Equal(ErrorLoadNotFound, err.(LoadErrors).GetCode())
}
