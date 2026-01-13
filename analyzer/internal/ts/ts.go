package ts

import (
	"fmt"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	tsjavascript "github.com/smacker/go-tree-sitter/javascript"
	tspython "github.com/smacker/go-tree-sitter/python"
	tstypescript "github.com/smacker/go-tree-sitter/typescript/typescript"
)

// language ids used by analyzer.
const (
	LangPython     = "python"
	LangJavaScript = "javascript"
	LangTypeScript = "typescript"
)

var languageMap = map[string]*sitter.Language{
	LangPython:     tspython.GetLanguage(),
	LangJavaScript: tsjavascript.GetLanguage(),
	LangTypeScript: tstypescript.GetLanguage(),
}

var extToLang = map[string]string{
	".py":  LangPython,
	".js":  LangJavaScript,
	".mjs": LangJavaScript,
	".cjs": LangJavaScript,
	".jsx": LangJavaScript,
	".ts":  LangTypeScript,
	".tsx": LangTypeScript,
}

// detect language from flag or extension.
func DetectLanguage(langFlag, path string) string {
	if langFlag != "" {
		return strings.ToLower(langFlag)
	}
	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := extToLang[ext]; ok {
		return lang
	}
	return ""
}

// parse source with language parser.
func Parse(lang string, source []byte) (*sitter.Node, error) {
	lang = strings.ToLower(lang)
	langImpl, ok := languageMap[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
	parser := sitter.NewParser()
	parser.SetLanguage(langImpl)
	tree := parser.Parse(nil, source)
	if tree == nil {
		return nil, fmt.Errorf("failed to parse source")
	}
	return tree.RootNode(), nil
}

// supported reports if language is available.
func Supported(lang string) bool {
	_, ok := languageMap[strings.ToLower(lang)]
	return ok
}
