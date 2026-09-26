package ir

type FileKind int

const (
	FileKindPage FileKind = iota
	FileKindComponent
	FileKindLayout
)

type GoImport struct {
	Alias string
	Path  string
}

type FileHeader struct {
	Kind      FileKind
	Name      string
	Props     []Prop
	Slots     []string
	Layout    string
	Profile   string
	Imports   []Import
	GoImports []GoImport
}

func (h FileHeader) IsComponent() bool {
	return h.Kind == FileKindComponent
}

func (h FileHeader) IsLayout() bool {
	return h.Kind == FileKindLayout
}
