package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"text/template"

	"github.com/iancoleman/strcase"
)

var target string
var cmdClass string

func init() {
	flag.StringVar(&cmdClass, "class", "", "command class")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatalf("target is required")
	}
	target = flag.Args()[0]
}

func main() {
	var defsFile string
	switch target {
	case "hostapi":
		defsFile = "ZWave_host_cmds.xml"
	case "commands":
		defsFile = "ZWave_cmd_classes.xml"
	}

	definitions, err := ioutil.ReadFile(fmt.Sprintf("gen/%s", defsFile))
	if err != nil {
		log.Fatal(err)
	}

	var doc Document
	err = xml.Unmarshal(definitions, &doc)
	if err != nil {
		log.Fatal(err)
	}
	temp := template.New("")
	temp, err = temp.Funcs(template.FuncMap{
		"toCamel":      strcase.ToCamel,
		"fieldName":    fieldName,
		"goTypeString": goTypeString,
	}).ParseGlob("gen/templates/*")
	if err != nil {
		log.Fatal(err)
	}
	m := new(sync.Mutex)
	sem := make(chan bool, 20)
	for _, cc := range doc.CommandClassDefs {
		if cmdClass != "" && cmdClass != cc.PackageName() {
			continue
		}
		sem <- true
		go generate(sem, m, temp, cc)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

type GetAndResponse struct {
	Get      *CommandDef
	Response *CommandDef
}

type GroupedCommandClassDef struct {
	CommandResponseDefs map[string]*GetAndResponse
	Others              []*CommandDef
}

func groupCommands(cc *CommandClassDef) {
	for i := 0; i < len(cc.CommandDefs); i++ {
		cmd := &cc.CommandDefs[i]
		cmd.Class = cc
		if len(cc.CommandDefs) > i+1 {
			neighbor := &cc.CommandDefs[i+1]
			neighbor.Class = cc
			if cmd.IsGet() && neighbor.IsReport() && cmd.ReportCommandName() == neighbor.ScreamingSnakeName {
				cmd.Report = neighbor
				i++
				continue
			}
		}
	}
	return
}

var controllerTemplate = `
package %s

import "encoding"

type Controller interface {
	SendAndReceive(encoding.BinaryMarshaler, encoding.BinaryUnmarshaler) error
	SendWithAcknowledgement(encoding.BinaryMarshaler) (int, error)
}
`

func generate(sem chan bool, m *sync.Mutex, temp *template.Template, cc CommandClassDef) {
	defer func() { <-sem }()
	ccDirectory := path.Join(target, cc.DirName())
	groupCommands(&cc)
	err := os.MkdirAll(ccDirectory, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	ctrlr := fmt.Sprintf(controllerTemplate, cc.PackageName())
	err = os.WriteFile(fmt.Sprintf("%s/%s.go", ccDirectory, "meta"), []byte(ctrlr), 0666)
	if err != nil {
		log.Fatal(err)
	}
	// for _, pair := range groups.CommandResponseDefs {
	// 	filePath := fmt.Sprintf("%s/%s.go", ccDirectory, pair.Get.FileName())
	// 	t := template.Must(template.New(filePath).Parse(string(temp)))
	// 	dest, err := os.Create(filePath)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer dest.Close()
	// 	err = t.Execute(dest, map[string]interface{}{"Command": pair.Get, "Response": pair.Response})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	for i := 0; i < len(cc.CommandDefs); i++ {
		filePath := fmt.Sprintf("%s/%s.go", ccDirectory, cc.CommandDefs[i].FileName())
		// t := template.Must(template.New(filePath).Parse(string(temp)))
		dest, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer dest.Close()
		err = temp.ExecuteTemplate(dest, "command.tpl", map[string]interface{}{"Command": &cc.CommandDefs[i]})
		if err != nil {
			log.Fatal(err)
		}
	}
	tree := exec.Command("tree", "--noreport", ccDirectory)
	tree.Stdout = os.Stdout
	m.Lock()
	err = tree.Run()
	if err != nil {
		log.Fatal(err)
	}
	m.Unlock()
}

func fieldName(param IParam) string {
	return strcase.ToCamel(invalidFieldChars.ReplaceAllString(param.Name(), ""))
}

func goTypeString(param IParam) string {
	switch param.Type() {
	case "ENUM":
		// we have declared a type
		return fieldName(param)
	case "BYTE":
		if param.ShowHex() {
			return "byte"
		} else {
			return "byte"
		}
	case "ARRAY":
		if param.ShowHex() {
			return "[]byte"
		} else {
			return "string"
		}
	case "WORD", "DWORD":
		return "byte"
	case "BIT_24":
		return "[]byte"
	case "CONST", "STRUCT_BYTE":
		return "byte"
	default:
		return param.Type()
	}
}
