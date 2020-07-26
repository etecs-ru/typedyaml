# typedyaml

This is a code generator for Go based on [github.com/etecs-ru/typedjson](https://github.com/etecs-ru/typedjson) that alleviates YAML marshaling/unmarshaling unrelated structs in typed fashion.

Imagine, that you send or receive a YAML object with some key `config`. 
The value of this field can correspond to two structs of your Go program: `FooConfig` and `BarConfig`.
So, the field `Config` in your struct must be able to hold a value of two possible types.
In this case, you have the following options:

1. You can declare field `Config` as `interface{}` and somehow determine what type you should expect, assign an object of this type to `Config` 
and then unmarshal object.
1. You can unmarshal field `Config` separately.
1. You can implement custom `MarshalYAML`/`UnmarshalYAML` for the third type that automatically will handle these cases.

This package provides means to generate all boilerplate code for the third case.

## Usage

```sh
typedyaml [OPTION] NAME...
```

Options:

* `-interface` string

	Name of the interface that encompass all types.

* `-output` string

	Output path where generated code should be saved.

* `-package` string

	Package name in generated file (default to GOPACKAGE).

* `-typed` string

	The name of the struct that will be used for typed interface (default to `{{interface}}{{Typed}}`).

Each name in position argument should be the name of the struct. 
You can set an alias for struct name like this: `foo=*FooConfig`.

## Example

For example, you have two structs:

```go
type FooConfig struct {
	Foo int
}

type BarConfig struct {
	Bar string
}
```

Then you must declare an interface that will hold either of these structs.
The interface must have the method `TypedYAML` with a special signature. 
This method will advise compiler to work with types.

```go
//go:generate go run github.com/etecs-ru/typedyaml -interface Config *FooConfig *BarConfig
type Config interface {
	TypedYAML(*ConfigTyped) string
}
```

After this, run `go generate`. 
Generated struct `ConfigTyped` will have special implemented methods `MarshalYAML`/`UnmarshalYAML`.
You can use generated code like this:

```go
type MyObject struct {
	Name string `yaml:"name"`
	ConfigTyped ConfigTyped `yaml:"config"`
}

func main() {
	yamlSources := []byte(`
        name: "my name"
        config:
            T: "*FooConfig",
            V:
                Foo: 123
    `)
	var object MyObject
	if err := yaml.Unmarshal(yamlSources, &object); err != nil {
		panic(err)
	}
	fooConfig := object.ConfigTyped.Config.(*FooConfig)
	if fooConfig.Foo != 123 {
		panic("fail")
	}
}
```
