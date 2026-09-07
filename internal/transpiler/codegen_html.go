package transpiler

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
	"github.com/dreego-stack/dreego/internal/transpiler/html/css"
	"github.com/dreego-stack/dreego/internal/transpiler/html/head"
	htmlinput "github.com/dreego-stack/dreego/internal/transpiler/html/html"
	"github.com/dreego-stack/dreego/internal/transpiler/html/output"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func genTemplateNode(gen *generator, n TemplateNode, depth int) (string, error) {
	return output.GenTemplateNode(gen, n, depth)
}

func genTemplateNodeToState(gen *generator, n TemplateNode, depth int, builder string, inSection *bool) (string, error) {
	return output.GenTemplateNodeToState(gen, n, depth, builder, inSection)
}

func genTempl(gen *Generator, file *File, layout *layoutEntry, scopeHash string, isGET bool) (string, error) {
	var l *codegen.Layout
	if layout != nil {
		l = &codegen.Layout{File: layout.file, Name: layout.name}
	}
	return output.GenTempl(gen, file, l, scopeHash, isGET)
}

func genTemplateNodeComp(gen *generator, n TemplateNode) (string, error) {
	return output.GenTemplateNodeComp(gen, n)
}

type compGen struct {
	gen     *generator
	builder string
}

func (g *compGen) genComponentCall(n TemplateNode) (string, error) {
	return output.GenComponentCall(g.gen, g.builder, n)
}

func GenerateComponent(gen *Generator, file *File, scopeHash string) (string, error) {
	return htmlinput.Generate(gen, file, scopeHash)
}

func buildComponentArgs(comp *ComponentDef, attrs string, src string, pos int) (string, error) {
	return output.BuildComponentArgs(comp, attrs, src, pos)
}

func componentParams(comp *ComponentDef) (decl, impl, call string, variadic string) {
	return htmlinput.Params(comp)
}

func writePropDefaultFallbacks(buf *strings.Builder, comp *ComponentDef) {
	htmlinput.WritePropDefaultFallbacks(buf, comp)
}

func validateSlotName(def *ComponentDef, name, filename, src string, pos int) error {
	return output.ValidateSlotName(def, name, filename, src, pos)
}

func nestedSlotError(call TemplateNode, def *ComponentDef, nested *TemplateNode, src string) error {
	return output.NestedSlotError(call, def, nested, src)
}

func genHead(htmlText string, bufName string) (string, error) {
	return head.Gen(htmlText, bufName)
}

func compTextWithAttrs(s string) string {
	return output.CompTextWithAttrs(s)
}

func compTextSection(content string, inSection bool) (string, bool) {
	return output.CompTextSection(content, inSection)
}

func attrSafeFunc(content string, tagStart, i int) string {
	return output.AttrSafeFunc(content, tagStart, i)
}

func scopeCSS(cssText string, hash string) string {
	return css.ScopeCSS(cssText, hash)
}

func scopeSelector(sel string, prefix string) string {
	return css.ScopeSelector(sel, prefix)
}

func splitTopLevelComma(sel string) []string {
	return css.SplitTopLevelComma(sel)
}

func matchBrace(cssText string, open, end int) int {
	return css.MatchBrace(cssText, open, end)
}

func headMergeHelpers() string {
	return output.HeadMergeHelpers()
}

func attrNameAt(tag string, i int) string {
	return ir.AttrNameAt(tag, i)
}

func attrContext(name string) string {
	return ir.AttrContext(name)
}

func isScriptAttr(name string) bool {
	return ir.IsScriptAttr(name)
}

func headSafeFunc(htmlText string, i int) string {
	return head.HeadSafeFunc(htmlText, i)
}

func attrValue(tag string, attr string) string {
	return ir.AttrValue(tag, attr)
}

func goLiteral(s string) string {
	return ir.GoLiteral(s)
}

func toPascalCase(s string) string {
	return ir.ToPascalCase(s)
}

func extractAttrValues(attrs string) string {
	return output.ExtractAttrValues(attrs)
}

func attrVal(part string) string {
	return output.AttrVal(part)
}

func concatPlaceholders(val string) string {
	return output.ConcatPlaceholders(val)
}
