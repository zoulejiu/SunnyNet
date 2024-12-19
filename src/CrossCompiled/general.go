package CrossCompiled

import (
	"github.com/qtgolang/SunnyNet/src/iphlpapi/net"
	"github.com/shirou/gopsutil/process"
	"os"
	"strconv"
)

// GetTcpInfoPID 用于获取指定 TCP 连接信息的 PID
func GetTcpInfoPID(tcpInfo string, SunnyPort int) string {
	connections, _ := net.Connections("tcp")
	for _, conn := range connections {
		if conn.Laddr.String() == tcpInfo {
			return strconv.Itoa(int(conn.Pid))
		}
	}
	return ""
}

// GetPidName 用于获取指定 PID 的进程名称
func GetPidName(pid int32) string {
	p, err := process.NewProcess(pid)
	if err != nil {
		return ""
	}
	name, err := p.Name()
	if err != nil {
		return ""
	}
	return name
}

var myPid = int32(os.Getpid())

// IsLoopRequest 是否环路请求
func IsLoopRequest(Port string, SunnyPort int) bool {
	p, _ := strconv.Atoi(Port)
	if p == 0 {
		return false
	}
	pp := uint32(p)
	pp2 := uint32(SunnyPort)
	connections, _ := net.ConnectionsPid("tcp", myPid)
	for _, conn := range connections {
		if conn.Laddr.Port == pp && conn.Raddr.Port == pp2 {
			return true
		}
	}
	return false
}

func LoopRemotePort(Srt string) uint32 {
	connections, _ := net.ConnectionsPid("tcp", myPid)
	for _, conn := range connections {
		if conn.Laddr.String() == Srt {
			return conn.Raddr.Port
		}
	}
	return 0
}
