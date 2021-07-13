package utils

import (
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/expr"
	"github.com/z7zmey/php-parser/node/expr/binary"
	"github.com/z7zmey/php-parser/node/stmt"
)

func GetAllClass(rootNode node.Node) {
	root := rootNode.(*node.Root)

	for _, s := range root.Stmts {
		if class, ok := s.(*stmt.Class); ok {
			clsName := class.ClassName.(*node.Identifier).Value
			if _, has := ClsList[clsName]; !has {
				ClsList[clsName] = make(map[string]*stmt.ClassMethod)
			}

			for _, m := range class.Stmts {
				if clsMethod, ok := m.(*stmt.ClassMethod); ok {
					methodName := clsMethod.MethodName.(*node.Identifier).Value
					ClsList[clsName][methodName] = clsMethod
				}
			}
		}
	}
}

func GetVarName(n node.Node) string {
	switch n.(type) {
	// this->xxx
	case *expr.PropertyFetch:
		varName := GetVarName(n.(*expr.PropertyFetch).Variable)
		property := GetVarName(n.(*expr.PropertyFetch).Property)
		return varName + "." + property
		// $xxx
	case *expr.Variable:
		return GetVarName(n.(*expr.Variable).VarName)
		// xxxx()
	case *node.Identifier:
		return n.(*node.Identifier).Value

	}
	return ""
}

// 题目中等号右边存在变量只有两种情况
func ResolveVarName(n node.Node) []string {
	switch n.(type) {
	case *expr.Variable:
		return []string{GetVarName(n)}
	case *binary.Concat:
		concat := n.(*binary.Concat)
		left := ResolveVarName(concat.Left)
		right := ResolveVarName(concat.Right)
		return append(left, right...)
	}

	return []string{}
}
