package transpiler

import "github.com/dreego-stack/dreego/internal/transpiler/ir"

type TemplateNodeType = ir.TemplateNodeType

const (
	NodeText          = ir.NodeText
	NodeExpression    = ir.NodeExpression
	NodeIf            = ir.NodeIf
	NodeEach          = ir.NodeEach
	NodeSlot          = ir.NodeSlot
	NodeComponentCall = ir.NodeComponentCall
	NodeVerbatim      = ir.NodeVerbatim
)

type TemplateNode = ir.TemplateNode
type ServerSection = ir.ServerSection
type File = ir.File
type ComponentDef = ir.ComponentDef
type Prop = ir.Prop
type Import = ir.Import
type HeadSection = ir.HeadSection
type BodySection = ir.BodySection
type ClientSection = ir.ClientSection
type StyleSection = ir.StyleSection
