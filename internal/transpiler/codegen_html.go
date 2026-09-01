package transpiler

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func genTemplateNode(gen *generator, n TemplateNode, depth int) (string, error) {
	return html.GenTemplateNode(gen, n, depth)
}

func genTemplateNodeToState(gen *generator, n TemplateNode, depth int, builder string, inSection *bool) (string, error) {
	return html.GenTemplateNodeToState(gen, n, depth, builder, inSection)
}

func genTempl(gen *Generator, file *File, layout *layoutEntry, scopeHash string, isGET bool) (string, error) {
	var l *ir.LayoutEntry
	if layout != nil {
		l = &ir.LayoutEntry{File: layout.file, Name: layout.name}
	}
	return html.GenTempl(gen, file, l, scopeHash, isGET)
}

func genTemplateNodeComp(gen *generator, n TemplateNode) (string, error) {
	return html.GenTemplateNodeComp(gen, n)
}

type compGen struct {
	gen     *generator
	builder string
}

func (g *compGen) genComponentCall(n TemplateNode) (string, error) {
	return html.GenComponentCall(g.gen, g.builder, n)
}

func GenerateComponent(gen *Generator, file *File, scopeHash string) (string, error) {
	return html.GenerateComponent(gen, file, scopeHash)
}

func buildComponentArgs(comp *ComponentDef, attrs string, src string, pos int) (string, error) {
	return html.BuildComponentArgs(comp, attrs, src, pos)
}

func componentParams(comp *ComponentDef) (decl, impl, call string, variadic string) {
	return html.ComponentParams(comp)
}

func writePropDefaultFallbacks(buf *strings.Builder, comp *ComponentDef) {
	html.WritePropDefaultFallbacks(buf, comp)
}

func validateSlotName(def *ComponentDef, name, filename, src string, pos int) error {
	return html.ValidateSlotName(def, name, filename, src, pos)
}

func nestedSlotError(call TemplateNode, def *ComponentDef, nested *TemplateNode, src string) error {
	return html.NestedSlotError(call, def, nested, src)
}

func genHead(htmlText string, bufName string) (string, error) {
	return html.GenHead(htmlText, bufName)
}

func compTextWithAttrs(s string) string {
	return html.CompTextWithAttrs(s)
}

func compTextSection(content string, inSection bool) (string, bool) {
	return html.CompTextSection(content, inSection)
}

func attrSafeFunc(content string, tagStart, i int) string {
	return html.AttrSafeFunc(content, tagStart, i)
}

func scopeCSS(css string, hash string) string {
	return html.ScopeCSS(css, hash)
}

func scopeSelector(sel string, prefix string) string {
	return html.ScopeSelector(sel, prefix)
}

func splitTopLevelComma(sel string) []string {
	return html.SplitTopLevelComma(sel)
}

func matchBrace(css string, open, end int) int {
	return html.MatchBrace(css, open, end)
}

func headMergeHelpers() string {
	return html.HeadMergeHelpers()
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
	return html.HeadSafeFunc(htmlText, i)
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
	return html.ExtractAttrValues(attrs)
}

func attrVal(part string) string {
	return html.AttrVal(part)
}

func concatPlaceholders(val string) string {
	return html.ConcatPlaceholders(val)
}
