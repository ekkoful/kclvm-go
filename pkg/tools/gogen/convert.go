package convert

import (
	"bytes"
	"fmt"
	"go/format"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

type GenStruct map[StructName][]*Field

type StructName string

type Field struct {
	Name       string
	SchemaType string
	OmitEmpty  bool
	Tag        string
}

func parseKclType(ktList []*gpyrpc.KclType) *GenStruct {
	var schemaName StructName
	genStruct := make(GenStruct, 0)
	for _, kt := range ktList {
		if kt.Type == "schema" {
			schemaName = StructName(kt.SchemaName)
			var schemaList []*Field
			for k, p := range kt.Properties {
				var sch Field
				sch.Name = k
				switch p.Type {
				case "schema":
					sch.SchemaType = "*" + p.SchemaName
				case "dict":
					kType := p.Key.Type
					if kType == "str" {
						kType = "string"
					} else if kType == "int" {
						kType = "int"
					}
					vType := p.Item.Type
					if vType == "schema" {
						vType = p.Item.SchemaName
					}
					sch.SchemaType = "map[" + kType + "]" + vType
				case "list":
					vType := p.Item.Type
					if vType == "schema" {
						vType = p.Item.SchemaName
					}
					sch.SchemaType = "[]" + "*" + vType
				case "str":
					sch.SchemaType = "string"
				case "int":
					sch.SchemaType = "int32"
				case "float":
					sch.SchemaType = "float32"
				case "bool":
					sch.SchemaType = "bool"
				case "null":
					sch.SchemaType = "nil"
				}
				schemaList = append(schemaList, &sch)
			}
			genStruct[schemaName] = schemaList
		}
	}
	return &genStruct
}

func GenGoCodeFromKclType(ktList []*gpyrpc.KclType) string {
	genStruct := parseKclType(ktList)
	var buf bytes.Buffer
	for k, v := range *genStruct {
		fmt.Fprintf(&buf, "type %s struct {\n", k)
		for _, field := range v {
			fmt.Fprintf(&buf, " %s %s \n", field.Name, field.SchemaType)
		}
		fmt.Fprintf(&buf, "}\n")
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Failed to format source:", err)
	}

	return string(source)
}
