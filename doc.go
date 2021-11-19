/*
Package config helps in adding configuration to
your application in an intutive way.

Example usage:

  // ...

  type MyConfiguration struct {
  	SomeValue string `default:"hello"`
  }

  func main() {
  	myConfig := MyConfiguration{}
  	if err := config.LoadConfiguration(&myConfig); err != nil {
			panic(err)
		}
		fmt.Println(myConfig.SomeValue)
  }

	// ...

In the above example, setting an environment
variable of `SOME_VALUE` will result in it being
consumed and printed
*/
package config
