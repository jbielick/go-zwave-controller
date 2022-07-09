// STOP
// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.

package api

const {{ .Name }} CommandID = {{ .ID }}

// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.
{{- range $type := .Types }}
//go:generate stringer -type={{ $type.Name }}
type {{ $type.Name }} {{ $type.Type }}
{{ end }}

// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.
const (
{{- range $const := .Constants }}
	{{ $const.Name }} {{ $const.Type }} = {{ $const.Value }}
{{- end }}
)

// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.
{{ if and .Response .Response.Fields -}}
type {{ $.Name }}Response struct {
	{{- range $field := .Response.Fields }}
	{{ $field.Name }} {{ $field.Type }}
	{{- end }}
}

func (r *{{ $.Name }}Response) UnmarshalBinary(data []byte) error {
	pos := 0
  {{- range $field := .Response.Fields }}
  {{- if $field.Length }}
  {{- if eq $field.Type "[]byte" }}
  r.{{ $field.Name }} = data[pos:pos+{{ $field.Length }}]
  {{- else }}
  r.{{ $field.Name }} = {{ $field.Type }}(data[pos:pos+{{ $field.Length }}])
  {{- end }}
  pos += {{ $field.Length }}
  {{- else }}
	r.{{ $field.Name }} = {{ if $field.Type }}{{ $field.Type }}{{ else }}{{ $field.Name }}{{ end }}(data[pos])
	pos++
  {{- end }}
  {{- end }}

	return nil
}
{{- end }}

{{ if eq .Response.Type "ack+res" -}}
func (c *Controller) {{ .Name }}() (*{{ $.Name }}Response, error) {
  r := &{{ $.Name }}Response{}
	frame, err := c.SendAndReceive(NewRequest({{ .Name }}, []byte{}))
  if err != nil {
    return r, err
  }
	r.UnmarshalBinary(frame.Payload)
	return r, nil
}
{{- else }}

{{- end }}
