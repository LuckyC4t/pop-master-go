package tracer

import (
	"fmt"
	"github.com/LuckyC4t/pop-master-go/internal/utils"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/expr"
	"github.com/z7zmey/php-parser/node/expr/assign"
	"github.com/z7zmey/php-parser/node/stmt"
	"os"
)

func trace(n node.Node, pollution map[string]bool) map[string]bool {
	switch n.(type) {
	case *assign.Assign:
		pollution = traceAssign(n, pollution)
	case *expr.Eval:
		pollution = traceEval(n, pollution)
	case *expr.MethodCall:
		pollution = traceMethodCall(n, pollution)
	case *stmt.Expression:
		pollution = traceExprssion(n, pollution)
	case *stmt.For:
		pollution = traceFor(n, pollution)
	case *stmt.If:
		pollution = traceIf(n, pollution)
	case *stmt.StmtList:
		pollution = traceStmtList(n, pollution)
	}

	return pollution
}

func traceEval(n node.Node, pollution map[string]bool) map[string]bool {
	eval := n.(*expr.Eval)
	arg := utils.GetVarName(eval.Expr)
	if pollution[arg] {
		for _, call := range utils.CallStack {
			fmt.Println("->", call)
		}
		// 只用找到一条有效路径就够了
		os.Exit(0)
	}
	return pollution
}

func traceAssign(n node.Node, pollution map[string]bool) map[string]bool {
	a := n.(*assign.Assign)
	leftName := utils.GetVarName(a.Variable)
	// right
	for _, name := range utils.ResolveVarName(a.Expression) {
		if pollution[name] {
			pollution[leftName] = true
			return pollution
		}
	}

	pollution[leftName] = false
	return pollution
}

func traceIf(n node.Node, pollution map[string]bool) map[string]bool {
	ifStmt := n.(*stmt.If)
	// condition
	pollution = trace(ifStmt.Cond, pollution)

	// stmt
	pollution = trace(ifStmt.Stmt, pollution)

	return pollution
}

func traceStmtList(n node.Node, pollution map[string]bool) map[string]bool {
	sl := n.(*stmt.StmtList)
	for _, st := range sl.Stmts {
		pollution = trace(st, pollution)
	}
	return pollution
}

func traceExprssion(n node.Node, pollution map[string]bool) map[string]bool {
	e := n.(*stmt.Expression)
	pollution = trace(e.Expr, pollution)
	return pollution
}

func traceFor(n node.Node, pollution map[string]bool) map[string]bool {
	forStmt := n.(*stmt.For)
	// init
	for _, e := range forStmt.Init {
		pollution = trace(e, pollution)
	}
	// cond
	for _, e := range forStmt.Cond {
		pollution = trace(e, pollution)
	}
	// loop
	for _, e := range forStmt.Loop {
		pollution = trace(e, pollution)
	}

	// stmts
	pollution = trace(forStmt.Stmt, pollution)

	return pollution
}

func traceMethodCall(n node.Node, pollution map[string]bool) map[string]bool {
	mc := n.(*expr.MethodCall)
	// arg
	arg := utils.GetVarName(mc.ArgumentList.Arguments[0].(*node.Argument).Expr)
	if pollution[arg] {
		// find target
		targetName := utils.GetVarName(mc.Method)
		for cls, methods := range utils.ClsList {
			for name, method := range methods {
				if name == targetName {
					// 找到目标，加入worklist
					clsM := utils.ClsMethod{
						Cls:    cls,
						Method: method,
					}
					utils.WorkList = append(utils.WorkList, clsM)
				}
			}
		}
	}

	return pollution
}
