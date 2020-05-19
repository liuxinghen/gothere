package gothere

// generate a field value from whole data map item,
// used in the scenario when an output field value relies on multiple input field value
type Generator func(data map[string]interface{}) (interface{}, error)

// convert an input field value to an output field value,
// used in the scenario when an output field value relies on one input field value
type Converter func(value interface{}) (interface{}, error)

// supplies a value as default value when required is false and input field value doesn't exist
type Supplier func() interface{}

type Rule struct {
	FromKey   string
	ToKey     string
	Required  bool
	Default   Supplier
	Converter Converter
	Generator Generator
	Mapping   map[interface{}]interface{}
}
