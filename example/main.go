package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	//"sigs.k8s.io/yaml"
)

type MyObject struct {
	Name        string      `yaml:"name"`
	ConfigTyped StructTyped `yaml:"config"`
}

type A struct {
	Foo string
	Bar int
}

func (a *A) TypedYAML(*StructTyped) string {
	return "a"
}

type B struct {
	Theta   float64
	Epsilon complex128
}

func (a *B) TypedYAML(*StructTyped) string {
	return "b"
}

type StructTyped struct {
	StructTypedConfig
}

type StructTypedConfig interface {
	TypedYAML(*StructTyped) string
}

func (t StructTyped) MarshalYAML() (interface{}, error) {
	if t.StructTypedConfig == nil {
		return nil, errors.New("nil interface in t.StructTypedConfig")
	}
	typedString := t.StructTypedConfig.TypedYAML(nil)
	wrapper := struct {
		T string
		V StructTypedConfig
	}{
		T: typedString,
		V: t.StructTypedConfig,
	}
	b, err := yaml.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *StructTyped) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var wrapperT struct {
		T string
	}
	if err := unmarshal(&wrapperT); err != nil {
		return err
	}

	switch wrapperT.T {
	case "a":
		var wrapperV struct {
			V *A
		}
		if err := unmarshal(&wrapperV); err != nil {
			return err
		}
		t.StructTypedConfig = wrapperV.V
	default:
		panic("oops")
	}
	return nil
}

func getCnf(typedString string) StructTypedConfig {
	switch typedString {
	case "a":
		return &A{
			Foo: "foo",
			Bar: 1,
		}

	case "b":
		return &B{
			Theta:   2.2,
			Epsilon: 0.11,
		}
	}
	return nil
}

func main() {
	type at struct {
		Name string `yaml:"name"`
		C    struct {
			T string
			V struct {
				Foo string
				Bar int
			}
		} `yaml:"config"`
	}

	a := at{
		Name: "my config",
		C: struct {
			T string
			V struct {
				Foo string
				Bar int
			}
		}{
			T: "foo",
			V: struct {
				Foo string
				Bar int
			}{
				Foo: "foo",
				Bar: 8,
			},
		},
	}

	abytes, err := yaml.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(abytes))

	yamlSources := []byte(`
name: my config
config:
  t: a
  v:
    foo: foo
    bar: 8
    `)
	var bb map[string]interface{}
	var object MyObject
	if err := yaml.Unmarshal(yamlSources, &bb); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlSources, &object); err != nil {
		panic(err)
	}
	fooConfig := object.ConfigTyped.StructTypedConfig.(*A)
	if fooConfig.Foo != "foo" {
		panic("fail")
	}

	b, err := yaml.Marshal(object)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
