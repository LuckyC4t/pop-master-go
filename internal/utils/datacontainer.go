package utils

import "github.com/z7zmey/php-parser/node/stmt"

type ClsMethod struct {
	Cls    string
	Method *stmt.ClassMethod
}

var ClsList map[string]map[string]*stmt.ClassMethod

var WorkList []ClsMethod

var CallStack []string

func (cm ClsMethod) String() string {
	return cm.Cls + "." + GetVarName(cm.Method.MethodName)
}
