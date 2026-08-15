package core

type TemplateNodeType int

const (
	NodeText TemplateNodeType = iota
	NodeExpression
	NodeIf
	NodeEach
	NodeSlot
	NodeComponentCall
	NodeVerbatim
)

type TemplateNode struct {
	Type         TemplateNodeType
	Content      string
	Cond         string
	Items        string
	Item         string
	Children     []TemplateNode
	ElseChildren []TemplateNode
	Tag          string
	Attrs        string
	SelfClose    bool
	Filters      []string
	Pos          int
	Source       string
}

type GoSection struct {
	Code           string
	Method         string
	MethodExplicit bool
	ContentType    string
	Action         string
}

type File struct {
	Head        *HeadSection
	Go          []GoSection
	Template    *TemplateSection
	Script      *ScriptSection
	Style       *StyleSection
	Component   *ComponentDef
	Imports     []Import
	FormActions []string
}

type ComponentDef struct {
	Name  string
	Props []Prop
	Slots []string
}

type Prop struct {
	Name    string
	Type    string
	Default string
}

type Import struct {
	Alias string
	Path  string
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
