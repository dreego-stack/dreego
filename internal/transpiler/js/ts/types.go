package ts

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func Declarations(component *ir.ComponentDef, sections []ir.ServerSection) (string, error) {
	var declarations strings.Builder
	if component != nil {
		fmt.Fprintf(&declarations, "interface %sProps {\n", component.Name)
		for _, prop := range component.Props {
			typeScriptType, err := goType(prop.Type)
			if err != nil {
				return "", fmt.Errorf("component %s prop %s: %w", component.Name, prop.Name, err)
			}
			fmt.Fprintf(&declarations, "  %s: %s;\n", prop.Name, typeScriptType)
		}
		declarations.WriteString("}\n")
	}
	for _, section := range sections {
		models, err := modelsFromGo(section.Code)
		if err != nil {
			return "", err
		}
		declarations.WriteString(models)
	}
	return declarations.String(), nil
}

func modelsFromGo(code string) (string, error) {
	file, _ := parser.ParseFile(token.NewFileSet(), "server.go", "package dreego\n"+code, parser.AllErrors)
	if file == nil {
		return "", nil
	}
	var out strings.Builder
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.TYPE {
			continue
		}
		for _, specification := range general.Specs {
			typeSpec := specification.(*ast.TypeSpec)
			structure, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			fmt.Fprintf(&out, "interface %s {\n", typeSpec.Name.Name)
			for _, field := range structure.Fields.List {
				fieldType, err := goASTType(field.Type)
				if err != nil {
					return "", fmt.Errorf("Go model %s: %w", typeSpec.Name.Name, err)
				}
				for _, name := range field.Names {
					if !name.IsExported() {
						continue
					}
					jsonName, include := modelFieldName(name.Name, field.Tag)
					if include {
						fmt.Fprintf(&out, "  %s: %s;\n", jsonName, fieldType)
					}
				}
			}
			out.WriteString("}\n")
		}
	}
	return out.String(), nil
}

func modelFieldName(fallback string, tag *ast.BasicLit) (string, bool) {
	if tag == nil {
		return fallback, true
	}
	value := strings.Trim(tag.Value, "`")
	jsonName, _, _ := strings.Cut(reflect.StructTag(value).Get("json"), ",")
	if jsonName == "-" {
		return "", false
	}
	if jsonName == "" {
		return fallback, true
	}
	return jsonName, true
}

func goASTType(expression ast.Expr) (string, error) {
	switch value := expression.(type) {
	case *ast.Ident:
		return goType(value.Name)
	case *ast.StarExpr:
		mapped, err := goASTType(value.X)
		return mapped + " | null", err
	case *ast.ArrayType:
		mapped, err := goASTType(value.Elt)
		return "Array<" + mapped + ">", err
	case *ast.MapType:
		key, err := goASTType(value.Key)
		if err != nil || key != "string" {
			return "", fmt.Errorf("only maps with string keys can be represented in TypeScript")
		}
		mapped, err := goASTType(value.Value)
		return "Record<string, " + mapped + ">", err
	case *ast.InterfaceType:
		return "unknown", nil
	case *ast.SelectorExpr:
		if pkg, ok := value.X.(*ast.Ident); ok && pkg.Name == "time" && value.Sel.Name == "Time" {
			return "string", nil
		}
		return value.Sel.Name, nil
	default:
		return "", fmt.Errorf("unsupported Go model field type %T", expression)
	}
}

func goType(value string) (string, error) {
	value = strings.TrimSpace(value)
	if after, ok := strings.CutPrefix(value, "*"); ok {
		mapped, err := goType(after)
		return mapped + " | null", err
	}
	if after, ok := strings.CutPrefix(value, "[]"); ok {
		mapped, err := goType(after)
		return "Array<" + mapped + ">", err
	}
	if after, ok := strings.CutPrefix(value, "map[string]"); ok {
		mapped, err := goType(after)
		return "Record<string, " + mapped + ">", err
	}
	switch value {
	case "string":
		return "string", nil
	case "bool":
		return "boolean", nil
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		return "number", nil
	case "any", "interface{}":
		return "unknown", nil
	default:
		if token.IsIdentifier(value) {
			return value, nil
		}
		return "", fmt.Errorf("unsupported Go type %q", value)
	}
}
