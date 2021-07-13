package parser

import (
	"errors"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/php7"
	"os"
)

func ParsePhpFile(src string) (node.Node, []error) {
	content, err := os.ReadFile(src)
	if err != nil {
		return nil, []error{err}
	}

	parser := php7.NewParser(content, "7.4")
	parser.Parse()

	parserErrs := parser.GetErrors()
	if len(parserErrs) != 0 {
		errs := make([]error, len(parserErrs))
		for i, e := range parserErrs {
			errs[i] = errors.New(e.String())
		}
		return nil, errs
	}

	rootNode := parser.GetRootNode()
	return rootNode, nil
}
