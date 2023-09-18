package autorest

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Package struct {
	Tag        string
	InputFiles []string
}

type YamlPackage struct {
	InputFiles []string `yaml:"input-file"`
}

var r = regexp.MustCompile(`\$\(tag\)\s+==\s+'(.+)'`)

func ParseAutoRestConfig(filename string) []Package {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	md := markdown.Parse(data, parser.NewWithExtensions(parser.NoExtensions))
	codeBlocks := allCodeBlocks(&md)

	out := make([]Package, 0)
	for _, codeBlock := range codeBlocks {
		if string(codeBlock.Info) == "yaml" {
			yamlPackage, err := ParseYamlConfig(string(codeBlock.Literal))
			if err != nil {
				logrus.Warnf("failed to parse yaml config: %+v", err)
			} else {
				for i, inputFile := range yamlPackage.InputFiles {
					yamlPackage.InputFiles[i] = path.Clean(path.Join(path.Dir(filename), inputFile))
				}
				out = append(out, *yamlPackage)
			}
		}
	}

	return out
}

func allCodeBlocks(node *ast.Node) []ast.CodeBlock {
	if node == nil {
		return nil
	}
	switch v := (*node).(type) {
	case *ast.Container:
		out := make([]ast.CodeBlock, 0)
		for _, child := range v.Children {
			out = append(out, allCodeBlocks(&child)...)
		}
		return out
	case *ast.Document:
		out := make([]ast.CodeBlock, 0)
		for _, child := range v.Children {
			out = append(out, allCodeBlocks(&child)...)
		}
		return out
	case *ast.Paragraph:
		out := make([]ast.CodeBlock, 0)
		for _, child := range v.Children {
			out = append(out, allCodeBlocks(&child)...)
		}
		return out
	case *ast.CodeBlock:
		return []ast.CodeBlock{*v}
	}
	return nil
}

func ParseYamlConfig(content string) (*Package, error) {
	matches := r.FindAllStringSubmatch(content, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		tag := matches[0][1]

		index := strings.Index(content, "\n")
		if index == -1 {
			return nil, fmt.Errorf("invalid yaml code block: no newline after tag, input: %v", content)
		}

		yamlContent := content[index+1:]
		var yamlPackage YamlPackage
		err := yaml.Unmarshal([]byte(yamlContent), &yamlPackage)
		if err != nil {
			return nil, err
		}

		if len(yamlPackage.InputFiles) == 0 {
			return nil, fmt.Errorf("input-file is empty, input: %v", content)
		}

		return &Package{
			Tag:        tag,
			InputFiles: yamlPackage.InputFiles,
		}, nil
	}
	return nil, fmt.Errorf("tag not found in yaml config: %s", content)
}
