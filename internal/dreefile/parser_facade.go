package dreefile

import (
	parserpkg "github.com/dreego-stack/dreego/internal/dreefile/parser"
	"github.com/dreego-stack/dreego/internal/dreefile/sections/body/md"
)

type Parser struct {
	inner *parserpkg.Parser
}

func NewParser(tokens []Token) *Parser {
	return &Parser{inner: parserpkg.NewParser(tokens)}
}

func NewParserConcatServer(tokens []Token) *Parser {
	return &Parser{inner: parserpkg.NewParserConcatServer(tokens)}
}

func (p *Parser) Parse() (*File, error) {
	file, err := p.inner.Parse()
	if err != nil {
		return nil, err
	}
	if file == nil {
		return file, nil
	}
	return file, processBodyMarkdown(file)
}

// processBodyMarkdown applies the Markdown transform once per distinct body and
// shares the resulting node slice between file.Body and its file.Bodies entry,
// so gogen.SetNodeSource still mutates the nodes the generator sees.
func processBodyMarkdown(file *File) error {
	processed := map[string][]TemplateNode{}
	var primaryKey string
	if file.Body != nil {
		nodes, err := md.ProcessBody(file.Body.Nodes, file.Body.Language)
		if err != nil {
			return err
		}
		file.Body.Nodes = nodes
		primaryKey = bodyKey(file.Body.Method, file.Body.Language)
	}
	for i := range file.Bodies {
		key := bodyKey(file.Bodies[i].Method, file.Bodies[i].Language)
		if key == primaryKey {
			file.Bodies[i].Nodes = file.Body.Nodes
			continue
		}
		if nodes, ok := processed[key]; ok {
			file.Bodies[i].Nodes = nodes
			continue
		}
		nodes, err := md.ProcessBody(file.Bodies[i].Nodes, file.Bodies[i].Language)
		if err != nil {
			return err
		}
		processed[key] = nodes
		file.Bodies[i].Nodes = nodes
	}
	return nil
}

func bodyKey(method, language string) string {
	return method + "\x00" + language
}
