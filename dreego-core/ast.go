package core

type TemplateNodeType int

const (
	NodeText TemplateNodeType = iota
	NodeExpression
	NodeIf
	NodeEach
	NodeSlot
)

type TemplateNode struct {
	Type     TemplateNodeType
	Content  string
	Cond     string
	Items    string
	Item     string
	Children []TemplateNode
}

type GoSection struct {
	Code   string
	Method string
}

type File struct {
	Head     *HeadSection
	Go       []GoSection
	Template *TemplateSection
	Script   *ScriptSection
	Style    *StyleSection
}

type HeadSection struct {
	Content string
}

type TemplateSection struct {
	Nodes []TemplateNode
}

type ScriptSection struct {
	Code string
}

type StyleSection struct {
	Code string
}
