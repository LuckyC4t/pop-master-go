package tracer

import (
	"fmt"
	"github.com/LuckyC4t/pop-master-go/internal/utils"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/stmt"
)

func Travel() {
	for len(utils.WorkList) > 0 {
		l := len(utils.WorkList)

		// pop
		current := utils.WorkList[0]
		utils.WorkList = utils.WorkList[1:]
		// push stack
		utils.CallStack = append(utils.CallStack, fmt.Sprint(current))

		// 函数内污点传播
		param := utils.GetVarName(current.Method.Params[0].(*node.Parameter).Variable)
		pollution := map[string]bool{
			param: true, // 如果能到达当前函数，说明参数可控
		}

		// trace stmts
		for _, st := range current.Method.Stmt.(*stmt.StmtList).Stmts {
			pollution = trace(st, pollution)
			if len(utils.WorkList) >= l {
				// 说明有新的函数
				Travel()
			}
		}

		// 当前函数分析完毕，离开
		utils.CallStack = utils.CallStack[:len(utils.CallStack)-1]
	}
}
