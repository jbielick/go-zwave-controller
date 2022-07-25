func (c *{{ .Command.StructName }}) UnmarshalBinary(data []byte) error {
{{- if ne (len .Command.AllParams) 0 }}
  {{- if .Command.Classless }}
  pos := 1 // skip command ID
  {{- else }}
  pos := 2 // skip class and command ID
  {{- end }}

  {{- range $param := .Command.AllParams }}
    {{- if eq $param.Type "ENUM" }}
  c.{{ fieldName . }} = {{ fieldName . }}(data[pos])
  pos++
    {{- else if eq $param.Type "ARRAY" }}
      {{- if $param.ArrayAttribute.ShowHex }}
  c.{{ fieldName $param }} = data[pos:pos+{{ $param.ArrayAttribute.Length }}]
      {{- else }}
  c.{{ fieldName $param }} = string(data[pos:pos+{{ $param.ArrayAttribute.Length }}])
      {{- end }}
  pos = pos+{{ $param.ArrayAttribute.Length }}
    {{- else if eq $param.Type "BYTE" }}
  c.{{ fieldName $param }} = data[pos]
  pos++
    {{- else if eq $param.Type "WORD" }}
  c.{{ fieldName $param }} = data[pos]
  pos++
    {{- else }}
  // marshal {{ $param.Key }} {{ $param.Index }}
  pos++
    {{- end }}
  {{- end }}
{{- end }}
  return nil
}
