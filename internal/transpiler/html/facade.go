package html

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html/css"
	"github.com/dreego-stack/dreego/internal/transpiler/html/head"
	"github.com/dreego-stack/dreego/internal/transpiler/html/html"
	"github.com/dreego-stack/dreego/internal/transpiler/html/output"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func GenTemplateNode(gen *ir.Generator, n ir.TemplateNode, depth int) (string, error) {
	return output.GenTemplateNode(gen, n, depth)
}

func GenTemplateNodeToState(gen *ir.Generator, n ir.TemplateNode, depth int, builder string, inSection *bool) (string, error) {
	return output.GenTemplateNodeToState(gen, n, depth, builder, inSection)
}

func GenTempl(gen *ir.Generator, file *ir.File, layout *ir.LayoutEntry, scopeHash string, isGET bool) (string, error) {
	return output.GenTempl(gen, file, layout, scopeHash, isGET)
}

func GenTemplateNodeComp(gen *ir.Generator, n ir.TemplateNode) (string, error) {
	return output.GenTemplateNodeComp(gen, n)
}

func GenComponentCall(gen *ir.Generator, builder string, n ir.TemplateNode) (string, error) {
	return output.GenComponentCall(gen, builder, n)
}

func CompTextWithAttrs(s string) string {
	return output.CompTextWithAttrs(s)
}

func CompTextSection(content string, inSection bool) (string, bool) {
	return output.CompTextSection(content, inSection)
}

func RestoreContextValue(indent, key, previous string) string {
	return output.RestoreContextValue(indent, key, previous)
}

func RestoreComponentContextValue(key, previous string) string {
	return output.RestoreComponentContextValue(key, previous)
}

func ComponentParams(comp *ir.ComponentDef) (decl, impl, call string, variadic string) {
	return html.ComponentParams(comp)
}

func WritePropDefaultFallbacks(buf *strings.Builder, comp *ir.ComponentDef) {
	html.WritePropDefaultFallbacks(buf, comp)
}

func ValidateSlotName(def *ir.ComponentDef, name, filename, src string, pos int) error {
	return output.ValidateSlotName(def, name, filename, src, pos)
}

func NestedSlotError(call ir.TemplateNode, def *ir.ComponentDef, nested *ir.TemplateNode, src string) error {
	return output.NestedSlotError(call, def, nested, src)
}

func FindNestedSlot(nodes []ir.TemplateNode) *ir.TemplateNode {
	return output.FindNestedSlot(nodes)
}

func BuildComponentArgs(comp *ir.ComponentDef, attrs string, src string, pos int) (string, error) {
	return output.BuildComponentArgs(comp, attrs, src, pos)
}

func GenerateComponent(gen *ir.Generator, file *ir.File, scopeHash string) (string, error) {
	return html.GenerateComponent(gen, file, scopeHash)
}

func GenHead(html string, bufName string) (string, error) {
	return head.Gen(html, bufName)
}

func ScopeCSS(cssText string, hash string) string {
	return css.ScopeCSS(cssText, hash)
}

func HeadMergeHelpers() string {
	return output.HeadMergeHelpers()
}

func SplitHeadPlaceholder(head string) (prefix, suffix string) {
	return output.SplitHeadPlaceholder(head)
}

func AttrSafeFunc(content string, tagStart, i int) string {
	return output.AttrSafeFunc(content, tagStart, i)
}

func HeadSafeFunc(html string, i int) string {
	return head.HeadSafeFunc(html, i)
}

func ExtractAttrValues(attrs string) string {
	return output.ExtractAttrValues(attrs)
}

func AttrVal(part string) string {
	return output.AttrVal(part)
}

func ConcatPlaceholders(val string) string {
	return output.ConcatPlaceholders(val)
}

func ScopeSelector(sel string, prefix string) string {
	return css.ScopeSelector(sel, prefix)
}

func SplitTopLevelComma(sel string) []string {
	return css.SplitTopLevelComma(sel)
}

func MatchBrace(cssText string, open, end int) int {
	return css.MatchBrace(cssText, open, end)
}
