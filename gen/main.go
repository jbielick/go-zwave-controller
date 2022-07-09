package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"text/template"

	"gopkg.in/yaml.v3"
)

// /*
// typedef unsigned char BYTE;
// typedef BYTE *        PBYTE;
// #include "ZW_classcmd.h"
// */
// import (
// 	"C"
// )

type data struct {
	Type string
	Name string
}

type Constant struct {
	Name  string
	Type  string
	Value string
}

type Type struct {
	Name string
	Type string
}

type ResponseField struct {
	Name   string
	Type   string
	Length int
}

type Response struct {
	Type   string
	Fields []ResponseField
}

type Command struct {
	Name      string
	ID        string
	Constants []Constant
	Types     []Type
	Response  Response
}

func main() {
	definitions, err := ioutil.ReadFile("gen/commands.yaml")

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]Command)
	// flag.StringVar(&d.Type, "type", "", "The subtype used for the queue being generated")
	// flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	// flag.Parse()
	err = yaml.Unmarshal(definitions, &data)
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	temp, err := ioutil.ReadFile("gen/command.go.tpl")
	if err != nil {
		log.Fatal(err)
	}
	for name, config := range data {
		wg.Add(1)
		config.Name = name
		go generate(wg, temp, config)
	}
	wg.Wait()
}

func generate(wg *sync.WaitGroup, temp []byte, cmd Command) {
	defer wg.Done()

	t := template.Must(template.New(cmd.Name).Parse(string(temp)))

	dest, err := os.Create(fmt.Sprintf("api/%s.go", cmd.Name))
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(dest, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
