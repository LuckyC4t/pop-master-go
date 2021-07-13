package main

import (
	"flag"
	"github.com/LuckyC4t/pop-master-go/internal/parser"
	"github.com/LuckyC4t/pop-master-go/internal/tracer"
	"github.com/LuckyC4t/pop-master-go/internal/utils"
	"github.com/z7zmey/php-parser/node/stmt"
	"log"
)

func main() {
	classFile := flag.String("file", "", "class file path")
	entryClass := flag.String("class", "", "entry class name")
	entryMethod := flag.String("method", "", "entry method name")

	flag.Parse()

	if *classFile == "" || *entryClass == "" || *entryMethod == "" {
		flag.Usage()
		return
	}

	rootNode, errs := parser.ParsePhpFile(*classFile)
	if len(errs) > 0 {
		for _, e := range errs {
			log.Println(e)
		}
		return
	}

	// init class list
	utils.ClsList = make(map[string]map[string]*stmt.ClassMethod)
	utils.GetAllClass(rootNode)

	entry := utils.ClsList[*entryClass][*entryMethod]
	utils.WorkList = append(utils.WorkList, utils.ClsMethod{
		Cls:    *entryClass,
		Method: entry,
	})

	// 进行遍历
	tracer.Travel()
}
