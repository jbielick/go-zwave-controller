func (c {{ .Command.StructName }}) MarshalBinary() ([]byte, error) {
  var payload []byte
  {{- if not .Command.Classless }}
  payload = append(payload, c.ClassID())
  {{- end }}
  payload = append(payload, c.ID())
  {{- range $param := .Command.AllParams }}
    {{- if eq $param.Type "ENUM" }}
  payload = append(payload, byte(c.{{ fieldName $param }}))
    {{- else if eq $param.Type "ARRAY" }}
      {{- if $param.ArrayAttribute.ShowHex }}
  payload = append(payload, c.{{ fieldName $param }}...)
      {{- else }}
  payload = append(payload, []byte(c.{{ fieldName $param }})...)
      {{- end }}
    {{- else if eq $param.Type "BYTE" }}
  payload = append(payload, c.{{ fieldName $param }})
    {{- else if eq $param.Type "WORD" }}
  payload = append(payload, c.{{ fieldName $param }})
    {{- else }}
  // marshal {{ $param.Key }} {{ $param.Index }}
    {{- end }}
  {{- end }}
  return payload, nil
}
