package utils

import (
	"github.com/astaxie/beego/logs"
	"runtime"
)

func InitLogsCores() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	runtime.GOMAXPROCS(runtime.NumCPU())
}
