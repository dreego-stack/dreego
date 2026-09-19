package dreefile

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/codegen"
	"github.com/dreego-stack/dreego/internal/dreefile/ir"
	bodyhtml "github.com/dreego-stack/dreego/internal/dreefile/sections/body/html"
	"github.com/dreego-stack/dreego/internal/dreefile/sections/head"
	"github.com/dreego-stack/dreego/internal/dreefile/sections/style"
)

func genTemplateNode(gen *generator, n TemplateNode, depth int) (string, error) {
	return bodyhtml.GenTemplateNode(gen, n, depth)
}

func genTemplateNodeToState(gen *generator, n TemplateNode, depth int, builder string, inSection *bool) (string, error) {
	return bodyhtml.GenTemplateNodeToState(gen, n, depth, builder, inSection)
}

func genTempl(gen *Generator, file *File, layout *layoutEntry, scopeHash string, isGET bool) (string, error) {
	var l *codegen.Layout
	if layout != nil {
		l = &codegen.Layout{File: layout.file, Name: layout.name}
	}
	return bodyhtml.GenTempl(gen, file, l, scopeHash, isGET)
}

func genTemplateNodeComp(gen *generator, n TemplateNode) (string, error) {
	return bodyhtml.GenTemplateNodeComp(gen, n)
}

type compGen struct {
	gen     *generator
	builder string
}

func (g *compGen) genComponentCall(n TemplateNode) (string, error) {
	return bodyhtml.GenComponentCall(g.gen, g.builder, n)
}

func GenerateComponent(gen *Generator, file *File, scopeHash string) (string, error) {
	return bodyhtml.Generate(gen, file, scopeHash)
}

func buildComponentArgs(comp *ComponentDef, attrs string, src string, pos int) (string, error) {
	return bodyhtml.BuildComponentArgs(comp, attrs, src, pos)
}

func componentParams(comp *ComponentDef) (decl, impl, call string, variadic string) {
	return bodyhtml.Params(comp)
}

func writePropDefaultFallbacks(buf *strings.Builder, comp *ComponentDef) {
	bodyhtml.WritePropDefaultFallbacks(buf, comp)
}

func validateSlotName(def *ComponentDef, name, filename, src string, pos int) error {
	return bodyhtml.ValidateSlotName(def, name, filename, src, pos)
}

func nestedSlotError(call TemplateNode, def *ComponentDef, nested *TemplateNode, src string) error {
	return bodyhtml.NestedSlotError(call, def, nested, src)
}

func genHead(htmlText string, bufName string) (string, error) {
	return head.Gen(htmlText, bufName)
}

func compTextWithAttrs(s string) string {
	return bodyhtml.CompTextWithAttrs(s)
}

func compTextSection(content string, inSection bool) (string, bool) {
	return bodyhtml.CompTextSection(content, inSection)
}

func attrSafeFunc(content string, tagStart, i int) string {
	return bodyhtml.AttrSafeFunc(content, tagStart, i)
}

func scopeCSS(cssText string, hash string) string {
	return style.ScopeCSS(cssText, hash)
}

func scopeSelector(sel string, prefix string) string {
	return style.ScopeSelector(sel, prefix)
}

func splitTopLevelComma(sel string) []string {
	return style.SplitTopLevelComma(sel)
}

func matchBrace(cssText string, open, end int) int {
	return style.MatchBrace(cssText, open, end)
}

func headMergeHelpers() string {
	return bodyhtml.HeadMergeHelpers()
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
