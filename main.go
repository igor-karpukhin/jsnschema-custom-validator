package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/igor-karpukhin/jsnschema-custom-validator/version"
	"github.com/pkg/errors"
)

const (
	ElID          = "$id"
	ElSchema      = "$schema"
	ElProperties  = "properties"
	ElType        = "type"
	ElDescription = "description"
	ElDefault     = "default"
	ElExamples    = "examples"
)

var (
	ErrNoID             = errors.New("no '$id' element found")
	ErrNoSchema         = errors.New("no '$schema' element found")
	ErrNoDescription    = errors.New("no 'description' element found")
	ErrNoRootProperties = errors.New("no 'properties' found in the root")
	ErrNoType           = errors.New("no 'type' element found")
	ErrNoDefault        = errors.New("no 'default' element found")
	ErrNoExamples       = errors.New("no 'examples' element found")
)

func die(msg string, err error) {
	if err != nil {
		fmt.Printf("%s. Error: %s\r\n", msg, err.Error())
	} else {
		fmt.Println(msg)
	}
	os.Exit(1)
}

func main() {
	oFile := flag.String("f", "", "JSON file name")
	oVersion := flag.Bool("v", false, "Prints the application version and exit")
	flag.Parse()

	if *oVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}
	if *oFile == "" {
		die("-f flag not provided or file name missing", nil)
	}

	fmt.Println("application started")
	fmt.Printf("trying to open file '%s'\r\n", *oFile)
	hFile, err := os.Open(*oFile)
	if err != nil {
		die("unable to open file "+*oFile, err)
	}

	var decoded map[string]interface{}
	err = json.NewDecoder(hFile).Decode(&decoded)
	if err != nil {
		die("unable to decode file "+*oFile, err)
	}

	err = validateRoot(decoded)
	if err != nil {
		die("unable to validate root element", err)
	}
	fmt.Println("root element is correct")

	errorsMap := make(map[string][]error)
	if props, ok := decoded[ElProperties].(map[string]interface{}); ok {
		validateProperties(props, errorsMap)
	} else {
		fmt.Println("incorrect root properties")
		os.Exit(-1)
	}

	if len(errorsMap) == 0 {
		fmt.Println("document is valid")
		os.Exit(0)
	}

	for eName, errs := range errorsMap {
		fmt.Println("Field: ", eName)
		for _, e := range errs {
			fmt.Println("\t-", e.Error())
		}
	}
	os.Exit(-1)
}

func validateRoot(m map[string]interface{}) error {
	if _, ok := m[ElID]; !ok {
		return ErrNoID
	}

	if _, ok := m[ElSchema]; !ok {
		return ErrNoSchema
	}

	if _, ok := m[ElProperties]; !ok {
		return ErrNoRootProperties
	}
	return nil
}

func validateProperties(m map[string]interface{},
	e map[string][]error) {
	for k, v := range m {
		if casted, ok := v.(map[string]interface{}); ok {
			if len(casted) > 0 {
				validateElement(k, casted, e)
			}
		}
	}
}

func validateElement(eName string, m map[string]interface{}, e map[string][]error) {
	propsFound := false
	for k, v := range m {
		if strings.ToLower(k) == ElProperties {
			if casted, ok := v.(map[string]interface{}); ok {
				validateProperties(casted, e)
			}
			propsFound = true
		}
	}
	if !propsFound {
		if _, ok := e[eName]; !ok {
			e[eName] = []error{}
		}
		if _, ok := m[ElDescription]; !ok {
			e[eName] = append(e[eName], ErrNoDescription)
		}
		if _, ok := m[ElDefault]; !ok {
			e[eName] = append(e[eName], ErrNoDefault)
		}
		if _, ok := m[ElExamples]; !ok {
			e[eName] = append(e[eName], ErrNoExamples)
		}
	}
	return
}
