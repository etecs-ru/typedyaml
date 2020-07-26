package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

const moduleName = "typedyaml"

const (
	ExitCodeError = 1
)

type Struct struct {
	Type  string
	Alias string
}

type GeneratorArgs struct {
	OutputPath string
	Interface  string
	Typed      string
	Package    string
	Imports    []string
	Structs    []Struct
	AllArgs    []string
}

func main() {
	args, err := parseArguments()
	if err != nil {
		exitf(ExitCodeError, "error while parsing arguments: %v\n", err)
	}
	buff := bytes.NewBuffer([]byte{})
	if err := generateCode(args, buff); err != nil {
		exitf(ExitCodeError, "error while generating code: %v\n", err)
	}

	code, err := imports.Process(filepath.Dir(args.OutputPath), buff.Bytes(), nil)
	if err != nil {
		exitf(ExitCodeError, "error while processing imports: %v\n", err)
	}

	if args.OutputPath == "stdout" {
		_, err = os.Stdout.Write(code)
	} else {
		err = ioutil.WriteFile(args.OutputPath, code, 0644)
	}
	if err != nil {
		exitf(ExitCodeError, "error while writing code to %s: %v\n", args.OutputPath, err)
	}
}

func parseStructs(args []string) []Struct {
	structs := make([]Struct, 0, len(args))
	for _, arg := range args {
		var s Struct
		if idx := strings.Index(arg, "="); idx == -1 {
			s.Alias = arg
			s.Type = arg
		} else {
			s.Alias = arg[:idx]
			s.Type = arg[idx+1:]
		}
		structs = append(structs, s)
	}
	return structs
}

func parseArguments() (*GeneratorArgs, error) {
	ga := GeneratorArgs{}
	flag.StringVar(&ga.Package, "package", os.Getenv("GOPACKAGE"), "package name in generated file (default to GOPACKAGE)")
	flag.StringVar(&ga.Interface, "interface", "", "name of the interface that encompass all types")
	flag.StringVar(&ga.Typed, "typed", "", "name of struct that will used for typed interface (default to %%interface%%Typed")
	flag.StringVar(&ga.OutputPath, "output", "", "output path where generated code should be saved")
	flag.Parse()

	ga.Structs = parseStructs(flag.Args())

	if ga.Typed == "" {
		ga.Typed = ga.Interface + "Typed"
	}

	if ga.OutputPath == "" {
		ga.OutputPath = strings.ToLower(fmt.Sprintf("%s_%s.go", ga.Interface, moduleName))
	}

	ga.AllArgs = os.Args
	ga.AllArgs[0] = moduleName
	if err := checkArgs(&ga); err != nil {
		return nil, err
	}
	return &ga, nil
}

func checkArgs(args *GeneratorArgs) error {
	if args.Package == "" {
		return errors.New("package name should not be empty")
	}
	if args.OutputPath == "" {
		return errors.New("output path should not be empty")
	}
	return nil
}

func exitf(code int, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...) // nolint: errcheck
	os.Exit(code)
}
