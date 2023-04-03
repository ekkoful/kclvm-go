package convert

import (
	"bytes"
	"fmt"
	"go/format"

	"kusionstack.io/kclvm-go/pkg/logger"
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
				var field Field
				field.Name = k
				switch p.Type {
				case "schema":
					field.SchemaType = "*" + p.SchemaName
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "schema")
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
					field.SchemaType = "map[" + kType + "]" + vType
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "{"+p.Key.Type+":"+vType+"}")
				case "list":
					vType := p.Item.Type
					if vType == "schema" {
						vType = p.Item.SchemaName
					}
					field.SchemaType = "[]" + "*" + vType
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "["+vType+"]")
				case "str":
					field.SchemaType = "string"
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "str")
				case "int":
					field.SchemaType = "int"
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "int")
				case "float":
					field.SchemaType = "float32"
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "float")
				case "bool":
					field.SchemaType = "bool"
					field.Tag = fmt.Sprintf(`kcl:"name:%s,type:%s"`, field.Name, "bool")
				case "null":
					field.SchemaType = "nil"
				}
				schemaList = append(schemaList, &field)
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
			if field.Tag != "" {
				fmt.Fprintf(&buf, " %s %s %s\n", field.Name, field.SchemaType, "`"+field.Tag+"`")
			} else {
				fmt.Fprintf(&buf, " %s %s\n", field.Name, field.SchemaType)
			}
		}
		fmt.Fprintf(&buf, "}\n")
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		logger.GetLogger().Errorf("Failed to format kclvm source code: %s", err.Error())
	}

	return string(source)
}
