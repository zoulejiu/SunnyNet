package SunnyNet

import (
	NFapi "SunnyNet/project/src/nfapi"
	"runtime"
)

func SetNetworkConnectNumber() {
	//https://blog.csdn.net/PYJcsdn/article/details/126251054
	//尽量避免这个问题
	if runtime.GOOS == "windows" {
		var args []string
		args = append(args, "int")
		args = append(args, "ipv4")
		args = append(args, "set")
		args = append(args, "dynamicport")
		args = append(args, "tcp")
		args = append(args, "start=10000")
		args = append(args, "num=55000")
		NFapi.ExecCommand("netsh", args)
		var args1 []string
		args1 = append(args1, "int")
		args1 = append(args1, "ipv6")
		args1 = append(args1, "set")
		args1 = append(args1, "dynamicport")
		args1 = append(args1, "tcp")
		args1 = append(args1, "start=10000")
		args1 = append(args1, "num=55000")
		NFapi.ExecCommand("netsh", args1)
	}
}
