package main

import "C"
import (
	_ "net/http/pprof"
	//不要问 问为什么要把http复制到项目里面来，不是多此一举？
	//只在这里解释一次，因为官方的http库 会自动添加UAgent，而且可能自动处理协议头大小写，那么复制到项目里面来，我们可以随意修改
)

func init() {
}

func main() {
	Test()

}
