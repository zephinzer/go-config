# go-config

An intuitive configuration loader for use in Golang applications.

- [go-config](#go-config)
- [Usage](#usage)
  - [Importing](#importing)
  - [Example](#example)
    - [A fully featured struct](#a-fully-featured-struct)
    - [Load()](#load)
- [Documentation](#documentation)
  - [Variable name parsing](#variable-name-parsing)
  - [Supported types of variables](#supported-types-of-variables)
  - [Supported struct tags](#supported-struct-tags)
- [Contributing](#contributing)
  - [Changelog](#changelog)
  - [Potential roadmap](#potential-roadmap)
- [License](#license)

# Usage

## Importing

```go
import "github.com/zephinzer/go-config"
```

## Example

### A fully featured struct

The following is a fully featured configuration `struct` which can be passed to the `.Load*()` methods for retrieving values from the environment.

```go
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
```

### Load()

See [./example/load/main.go](./example/load/main.go) for a fully-featured example.

To get an idea of how the different struct tags work, clone this repository and try the example by running:

```sh
CUSTOM_BOOL=true \
  REQUIRED_BOOL=true \
  REQUIRED_INT=-54321 \
  REQUIRED_STRING="hello world" \
  REQUIRED_STRING_SLICE="hello,world" \
  go run ./example/load;
```

You should get the following pretty-printed JSON output:

```json
{
  "OptionalBool": null,
  "OptionalBoolWithCustomEnv": true,
  "OptionalInt": null,
  "OptionalString": null,
  "OptionalStringSlice": null,
  "RequiredBool": true,
  "RequiredBoolWithDefault": true,
  "RequiredInt": -54321,
  "RequiredIntWithDefault": -12345,
  "RequiredString": "hello world",
  "RequiredStringWithDefault": "required-string",
  "RequiredStringSlice": [
    "hello",
    "world"
  ],
  "RequiredStringSliceWithDefault": [
    "required",
    "string",
    "slice"
  ],
  "RequiredStringSliceWithDefaultAndDelimiter": [
    "required",
    "string",
    "slice"
  ]
}
```

# Documentation

## Variable name parsing

This library converts variable names in your provided configuration `struct` into `UPPER_SNAKE_CASE` and checks the environment for these keys. The library `go-strcase` is used to parse the variable names into environment variables. Here are some examples of the transformation:

| Variable name | Derived environment value key |
| --- | --- |
| someString | `SOME_STRING` |
| SomeNumber | `SOME_NUMBER` |
| endpointURL | `ENDPOINT_URL` |
| IsURLSetCorrectly | `IS_URL_SET_CORRECTLY` |

## Supported types of variables

| Type | Golang type |
| --- | --- |
| Boolean | `bool` |
| Optional boolean | `*bool` |
| Number | `int` |
| Optional number | `*int` |
| String | `string` |
| Optional string | `*string` |
| String Slice | `[]string` |
| Optional string slice | `*[]string` |

## Supported struct tags

Struct tags are used to provide metadata for the `.Load*()` methods to process. The following are the available struct tags:

| Struct tag | Example | Description |
| --- | --- | --- |
| `default` | `default:"153"` | Defines a default value for the variable. This has to be specified as a string; if the type of the property is not a string, the library parses the provided string value into the required type. |
| `delimiter` | `delimiter:"custom_bool"` | Applies only to slice types and defines the delimiter which should be used. This defaults to the comma character - `,` |
| `env` | `env:"custom_bool"` | Defines the environment variable key to use instead of the default one derived from the property name |

# Contributing

The working repository is at Gitlab at [https://gitlab.com/zephinzer/go-config](https://gitlab.com/zephinzer/go-config), if you are seeing this on Github, it's just for SEO since y'know all the cool new kids are on Github 😂

## Changelog

A rough changelog when the contributors can remember to add it is here:

| Version  | Description                                                                                                                                 |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------- |                                         |
| `v1.0.0` | Initial release |

## Potential roadmap

If you'd like to contribute, here are some features I have a distant-future need to implement but haven't gotten round to doing:

- [ ] Implement a `.Load*()` method that consumes variable values from derived flag names (probably convert to `--lower-kebab-case` and provide a struct tag to define the shorthand flag notation)
- [ ] Implement a `.Load*()` method that consumes variable values from a file instead of environment variables
- [ ] Integration with `spf13/cobra` via `PreRun` or `PersistentPreRun` to automatically run and get values
- [ ] Integration with `spf13/cobra` AND `spf13/viper` via `PersistentFlags` to automatically run and get values from flags

# License

Use this anywhere you need to. Licensed under [the MIT license](./LICENSE)
