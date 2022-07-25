package main

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type BasicDeviceDef struct {
	XMLName  xml.Name `xml:"bas_dev"`
	ReadOnly bool     `xml:"read_only,attr"`
	Name     string   `xml:"name,attr"`
	Key      string   `xml:"key,attr"`
	Help     string   `xml:"help,attr"`
	Comment  string   `xml:"comment,attr"`
}

type GenericDeviceDef struct {
	BasicDeviceDef
	XMLName            xml.Name            `xml:"gen_dev"`
	SpecificDeviceDefs []SpecificDeviceDef `xml:"spec_dev"`
}

type SpecificDeviceDef struct {
	BasicDeviceDef
	XMLName xml.Name `xml:"spec_dev"`
}

type CommandClassDef struct {
	XMLName            xml.Name     `xml:"cmd_class"`
	Key                string       `xml:"key,attr"`
	Version            string       `xml:"version,attr"`
	ScreamingSnakeName string       `xml:"name,attr"`
	Help               string       `xml:"help,attr"`
	Comment            string       `xml:"comment,attr"`
	ReadOnly           bool         `xml:"read_only,attr"`
	CommandDefs        []CommandDef `xml:"cmd"`
}

func (cc *CommandClassDef) UnprefixedName() string {
	return strings.Replace(cc.ScreamingSnakeName, "COMMAND_CLASS_", "", 1)
}

func (cc *CommandClassDef) UnprefixedSnakeName() string {
	return strcase.ToSnake(cc.UnprefixedName())
}

func (cc *CommandClassDef) UnprefixedCamelName() string {
	return strcase.ToCamel(cc.UnprefixedName())
}

func (cc *CommandClassDef) PackageName() string {
	return strings.ToLower(
		strings.Replace(
			cc.UnprefixedName(),
			"_", "", -1,
		),
	)
}

func (cc *CommandClassDef) DirName() string {
	return path.Join(cc.PackageName(), fmt.Sprintf("v%s", cc.Version))
}

type CommandDef struct {
	XMLName            xml.Name `xml:"cmd"`
	Key                string   `xml:"key,attr"`
	ScreamingSnakeName string   `xml:"name,attr"`
	ShortenedName      string
	Help               string `xml:"help,attr"`
	Comment            string `xml:"comment,attr"`
	Class              *CommandClassDef
	Params             []CommandDefParam `xml:"param"`
	VariantGroups      []VariantGroup    `xml:"variant_group"`
	Report             *CommandDef
}

func (c *CommandDef) NonRedundantName() string {
	return strings.Replace(
		c.ScreamingSnakeName,
		fmt.Sprintf("%s_", c.Class.UnprefixedName()),
		"",
		1,
	)
}

func (c *CommandDef) StructName() string {
	return strcase.ToCamel(c.NonRedundantSnakeName())
}

func (c *CommandDef) NonRedundantSnakeName() string {
	return strcase.ToSnake(c.NonRedundantName())
}

func (c *CommandDef) FileName() string {
	return c.NonRedundantSnakeName()
}

func (c *CommandDef) IsGet() bool {
	return strings.HasSuffix(c.ScreamingSnakeName, "_GET")
}

func (c *CommandDef) ReportCommandName() string {
	return fmt.Sprintf(
		"%s_REPORT",
		strings.Replace(c.ScreamingSnakeName, "_GET", "", 1),
	)
}

func (c *CommandDef) IsReport() bool {
	return strings.HasSuffix(c.ScreamingSnakeName, "_REPORT")
}

func (c *CommandDef) GetCommandName() string {
	return fmt.Sprintf(
		"%s_GET",
		strings.Replace(c.ScreamingSnakeName, "_REPORT", "", 1),
	)
}

func (c *CommandDef) Classless() bool {
	return c.Class.Key == "0x00"
}

type VariantGroup struct {
	XMLName      xml.Name          `xml:"variant_group"`
	Key          string            `xml:"key,attr"`
	name         string            `xml:"name,attr"`
	VariantKey   string            `xml:"variantKey,attr"`
	ParamOffs    string            `xml:"paramOffs,attr"`
	SizeMask     string            `xml:"sizemask,attr"`
	SizeOffs     string            `xml:"sizeoffs,attr"`
	typeHashCode string            `xml:"typehashcode,attr"`
	Params       []CommandDefParam `xml:"param"`
}

func (g *VariantGroup) Type() string {
	return "VG"
}

func (g *VariantGroup) Name() string {
	return g.name
}

func (g *VariantGroup) TypeHashCode() string {
	return g.typeHashCode
}

func (p *VariantGroup) Index() byte {
	b, _ := hex.DecodeString(p.Key[2:])
	return b[0]
}

func (g *VariantGroup) ShowHex() bool {
	return true
}

type EnumValue struct {
	Key  string `xml:"key,attr"`
	Name string `xml:"name,attr"`
}

type ArrayAttribute struct {
	Key     string `xml:"key,attr"`
	Length  int    `xml:"len,attr"`
	IsASCII bool   `xml:"is_ascii,attr"`
	ShowHex bool   `xml:"showhex,attr"`
}

type CommandDefParam struct {
	XMLName        xml.Name                       `xml:"param"`
	Key            string                         `xml:"key,attr"`
	ParamName      string                         `xml:"name,attr"`
	ParamType      string                         `xml:"type,attr"`
	ParamHashCode  string                         `xml:"typehashcode,attr"`
	Comment        string                         `xml:"comment,attr"`
	ValueAttribute *CommandDefParamValueAttribute `xml:"valueattrib"`
	Constants      []CommandDefParamConstant      `xml:"const"`
	BitMask        []CommandDefParamBitMask       `xml:"bitmask"`
	BitField       []CommandDefParamBitField      `xml:"bitfield"`
	ArrayAttribute *ArrayAttribute                `xml:"arrayattrib"`
	// <arrayattrib key="0x00" len="16" is_ascii="false" showhex="true" />
	// FieldEnum []
	// <fieldenum key="0x00" fieldname="Pay Meter" fieldmask="0x0F" shifter="0">
	// <fieldenum value="Reserved" />
	// <fieldenum value="Creditmeter" />
	// <fieldenum value="Prepayment meter" />
	// <fieldenum value="Prepayment meter debt" />
	// </fieldenum>
	// Variant []
	// <variant paramoffs="0" showhex="true" signed="true" sizemask="0x1F" sizeoffs="0" />
	// Bit24 []
	// <bit_24 key="0x00" hasdefines="false" showhex="true" />
	// Word
	// <word key="0x00" hasdefines="false" showhex="true" />
	// DWord
	// <dword key="0x00" hasdefines="false" showhex="true" />
	EnumValues []EnumValue `xml:"enum"`
	// <enum key="0x71" name="COMMAND_CLASS_ALARM" />
	// MultiArray
	// <multi_array>
	// <paramdescloc key="0x00" param="4" paramdesc="255" paramstart="4" />
	// </multi_array>
	// <multi_array>
	// <bitflag key="0x04" flagname="SPECIFIC_TYPE_NOT_USED" flagmask="0x00" />
	// <bitflag key="0x04" flagname="SPECIFIC_TYPE_DOORBELL" flagmask="0x12" />
	// <bitflag key="0x04" flagname="SPECIFIC_TYPE_SATELLITE_RECEIVER" flagmask="0x04" />
	// <bitflag key="0x04" flagname="SPECIFIC_TYPE_SATELLITE_RECEIVER_V2" flagmask="0x11" />
	// <bitflag key="0x04" flagname="SPECIFIC_TYPE_SOUND_SWITCH" flagmask="0x01" />
	// </multi_array>
}

var invalidFieldChars = regexp.MustCompile(`[^a-zA-Z0-9]`)

func (p *CommandDefParam) Type() string {
	return p.ParamType
}

func (p *CommandDefParam) Name() string {
	return p.ParamName
}

func (p *CommandDefParam) Index() byte {
	b, _ := hex.DecodeString(p.Key[2:])
	return b[0]
}

func (p *CommandDefParam) TypeHashCode() string {
	return p.ParamHashCode
}

func (p *CommandDefParam) ShowHex() bool {
	if p.ValueAttribute != nil {
		return p.ValueAttribute.ShowHex
	} else if p.ArrayAttribute != nil {
		return p.ArrayAttribute.ShowHex
	} else {
		return false
	}
}

type IParam interface {
	Type() string
	Name() string
	ShowHex() bool
	Index() byte
	TypeHashCode() string
}

func (c *CommandDef) AllParams() []IParam {
	all := make([]IParam, len(c.Params)+len(c.VariantGroups))
	for i := 0; i < len(c.Params); i++ {
		param := &c.Params[i]
		all[param.Index()] = param
	}
	for i := 0; i < len(c.VariantGroups); i++ {
		group := &c.VariantGroups[i]
		all[group.Index()] = group
	}
	return all
}

type CommandDefParamValueAttribute struct {
	XMLName    xml.Name `xml:"valueattrib"`
	Key        string   `xml:"key,attr"`
	HasDefines bool     `xml:"hasdefines,attr"`
	ShowHex    bool     `xml:"showhex,attr"`
}

type CommandDefParamConstant struct {
	XMLName  xml.Name `xml:"const"`
	Key      string   `xml:"key,attr"`
	FlagName string   `xml:"flagname,attr"`
	FlagMask string   `xml:"flagmask,attr"`
}

type CommandDefParamBitMask struct {
	XMLName     xml.Name `xml:"bitmask"`
	Key         string   `xml:"key,attr"`
	ParamOffset int      `xml:"paramoffs,attr"`
	LenMask     string   `xml:"lenmask,attr"`
	LenOffset   int      `xml:"lenoffs,attr"`
	Len         int      `xml:"len,attr"`
}

type CommandDefParamBitField struct {
	XMLName   xml.Name `xml:"bitfield"`
	Key       string   `xml:"key,attr"`
	FieldName string   `xml:"fieldname,attr"`
	FieldMask string   `xml:"fieldmask,attr"`
	Shifter   int      `xml:"shifter,attr"`
}

type Document struct {
	XMLName           xml.Name           `xml:"zw_classes"`
	BasicDeviceDefs   []BasicDeviceDef   `xml:"bas_dev"`
	GenericDeviceDefs []GenericDeviceDef `xml:"gen_dev"`
	CommandClassDefs  []CommandClassDef  `xml:"cmd_class"`
}

// func longestCommonPrefix(strs []string) string {
// 	var longestPrefix string = ""
// 	var endPrefix = false

// 	if len(strs) > 0 {
// 		sort.Strings(strs)
// 		first := string(strs[0])
// 		last := string(strs[len(strs)-1])

// 		for i := 0; i < len(first); i++ {
// 			if !endPrefix && string(last[i]) == string(first[i]) {
// 				longestPrefix += string(last[i])
// 			} else {
// 				endPrefix = true
// 			}
// 		}
// 	}
// 	return longestPrefix
// }
