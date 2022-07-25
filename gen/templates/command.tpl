// STOP
// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.

package {{ .Command.Class.PackageName }} // {{ .Command.Class.Key }}

{{- range $param := .Command.AllParams }}
{{- if eq $param.Type "ENUM" }}
//go:generate stringer -type={{ fieldName $param }}
type {{ fieldName $param }} byte

const (
{{- range $e := $param.EnumValues }}
  {{ toCamel $e.Name }} {{ fieldName $param }} = {{ $e.Key }}
{{- end }}
)
{{- end }}
{{- end }}

type {{ .Command.StructName }} struct {
  {{- range $param := .Command.AllParams }}
  {{ fieldName $param }} {{ goTypeString $param }} // {{ $param.Key }} {#{ with $param.Comment }}{#{ . }}{#{ end }}
  {{- end }}
}

func New{{ .Command.StructName }}() {{ .Command.StructName }} {
  return {{ .Command.StructName }}{}
}

func (c {{ $.Command.StructName }}) ClassID() byte {
  return {{ .Command.Class.Key }}
}

func (c {{ .Command.StructName }}) ID() byte {
  return {{ .Command.Key }}
}

func (c {{ .Command.StructName }}) Name() string {
  return "{{ .Command.ScreamingSnakeName }}"
}

func (c {{ .Command.StructName }}) Help() string {
  return "{{ .Command.Help }}"
}

func (c {{ .Command.StructName }}) Comment() string {
  return "{{ .Command.Help }}"
}

{{ template "unmarshal_binary.tpl" . }}

{{ template "marshal_binary.tpl" . }}

{{- if and .Command.IsGet .Command.Report }}
func (cmd {{ .Command.StructName }}) Send(c Controller) ({{ .Command.Report.StructName }}, error) {
	r := {{ .Command.Report.StructName }}{}
	err := c.SendAndReceive(cmd, &r)
	return r, err
{{- else }}
func (cmd *{{ .Command.StructName }}) Send(c Controller) (error) {
  _, err := c.SendWithAcknowledgement(cmd)
  return err
{{- end }}
}
