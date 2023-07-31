package iphlpapi

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"strconv"
	"syscall"
	"unsafe"
)

type (
	TCP_TABLE_CLASS int32
	DWORD           uint32
	ULONG           uint32
	MIB_TCP_STATE   int32
)

type MIB_TCPROW2 struct {
	dwState      uint32
	dwLocalAddr  uint32
	dwLocalPort  uint32
	dwRemoteAddr uint32
	dwRemotePort uint32
	dwOwningPid  uint32
}

type MIB_TCPTABLE2 struct {
	dwNumEntries uint32
	table        [1]MIB_TCPROW2
}

const (
	// TCP连接表类型常量
	TCP_TABLE_BASIC_LISTENER           TCP_TABLE_CLASS = iota
	TCP_TABLE_BASIC_CONNECTIONS                        // 基本连接
	TCP_TABLE_BASIC_ALL                                // 基本所有
	TCP_TABLE_OWNER_PID_LISTENER                       // PID监听
	TCP_TABLE_OWNER_PID_CONNECTIONS                    // PID连接
	TCP_TABLE_OWNER_PID_ALL                            // PID所有
	TCP_TABLE_OWNER_MODULE_LISTENER                    // 模块监听
	TCP_TABLE_OWNER_MODULE_CONNECTIONS                 // 模块连接
	TCP_TABLE_OWNER_MODULE_ALL                         // 模块所有
)

const (
	// MIB_TCP_STATE常量
	MIB_TCP_STATE_CLOSED     MIB_TCP_STATE = 1
	MIB_TCP_STATE_LISTEN     MIB_TCP_STATE = 2
	MIB_TCP_STATE_SYN_SENT   MIB_TCP_STATE = 3
	MIB_TCP_STATE_SYN_RCVD   MIB_TCP_STATE = 4
	MIB_TCP_STATE_ESTAB      MIB_TCP_STATE = 5
	MIB_TCP_STATE_FIN_WAIT1  MIB_TCP_STATE = 6
	MIB_TCP_STATE_FIN_WAIT2  MIB_TCP_STATE = 7
	MIB_TCP_STATE_CLOSE_WAIT MIB_TCP_STATE = 8
	MIB_TCP_STATE_CLOSING    MIB_TCP_STATE = 9
	MIB_TCP_STATE_LAST_ACK   MIB_TCP_STATE = 10
	MIB_TCP_STATE_TIME_WAIT  MIB_TCP_STATE = 11
	MIB_TCP_STATE_DELETE_TCB MIB_TCP_STATE = 12
)

var (
	// 加载iphlpapi.dll库
	lib                  = syscall.MustLoadDLL("iphlpapi.dll")
	GetTcpTable2         = lib.MustFindProc("GetTcpTable2").Addr()
	_setTCPEntry         = lib.MustFindProc("SetTcpEntry").Addr()
	_getExtendedTCPTable = lib.MustFindProc("GetExtendedTcpTable").Addr()
	CloseSocketErr       = errors.New("error closing game socket")
)

// setTCPEntry 设置TCP连接的状态
func setTCPEntry(ptr uintptr) error {
	_, _, err := syscall.SyscallN(_setTCPEntry, ptr)
	return err
}

// GetExtendedTCPTable 获取扩展TCP连接表
func getExtendedTCPTable(tcpTablePtr uintptr, pdwSize *DWORD, bOrder bool, ulAf ULONG, tableClass TCP_TABLE_CLASS, reserved ULONG) error {
	_, _, err := syscall.SyscallN(_getExtendedTCPTable,
		tcpTablePtr,
		uintptr(unsafe.Pointer(pdwSize)),
		getUintptrFromBool(bOrder),
		uintptr(ulAf),
		uintptr(tableClass),
		uintptr(reserved))

	return err
}

// getUintptrFromBool 将bool类型转换为uintptr类型
func getUintptrFromBool(b bool) uintptr {
	if b {
		return 1
	} else {
		return 0
	}
}

// CloseCurrentSocket 关闭指定地址族下的指定进程的所有TCP连接
func CloseCurrentSocket(PID int, ulAf uint) error {
	var lastOpenSocket []byte
	var buffSize DWORD
	_ = getExtendedTCPTable(uintptr(0), &buffSize, true, ULONG(ulAf), TCP_TABLE_OWNER_PID_ALL, 0)
	var buffTable = make([]byte, int(buffSize))
	err := getExtendedTCPTable(uintptr(unsafe.Pointer(&buffTable[0])), &buffSize, true, ULONG(ulAf), TCP_TABLE_OWNER_PID_ALL, 0)
	if !errors.Is(err, windows.DNS_ERROR_RCODE_NO_ERROR) {
		return CloseSocketErr
	}
	count := *(*uint32)(unsafe.Pointer(&buffTable[0]))
	const structLen = 24
	for n, pos := uint32(0), 4; n < count && pos+structLen <= len(buffTable); n, pos = n+1, pos+structLen {
		state := *(*uint32)(unsafe.Pointer(&buffTable[pos]))
		if state < 1 || state > 12 {
			return CloseSocketErr
		}
		pid := *(*uint32)(unsafe.Pointer(&buffTable[pos+20]))
		if PID == -1 || (uint(pid) == uint(PID) && state == uint32(MIB_TCP_STATE_ESTAB)) { // 判断PID是否匹配，并且连接状态是否为已建立
			buffTable[pos] = byte(MIB_TCP_STATE_DELETE_TCB) // 将连接状态设置为MIB_TCP_STATE_DELETE_TCB，表示要关闭该连接
			lastOpenSocket = buffTable[pos : pos+24]        // 记录最后一个打开的套接字
			if len(lastOpenSocket) == 0 {
				return CloseSocketErr
			}
			err = setTCPEntry(uintptr(unsafe.Pointer(&lastOpenSocket[0]))) // 设置连接状态为MIB_TCP_STATE_DELETE_TCB，关闭连接
		}
	}
	if errors.Is(err, windows.DNS_ERROR_RCODE_NO_ERROR) {
		return nil
	}
	return CloseSocketErr
}

// GetTcpInfoPID 用于获取指定 TCP 连接信息的 PID
func GetTcpInfoPID(tcpInfo string, SunnyPort int) string {
	var tcpTable *MIB_TCPTABLE2 // 声明 MIB_TCPTABLE2 结构体指针，用于保存 TCP 连接信息
	var size uint32             // 声明 size 变量，用于保存 TCP 连接信息的大小

	// 调用 syscall.Syscall 函数调用 GetTcpTable2 函数，获取 TCP 连接信息的大小
	ret, _, _ := syscall.Syscall(GetTcpTable2, 3, uintptr(unsafe.Pointer(tcpTable)), uintptr(unsafe.Pointer(&size)), 1)
	if ret != 122 { // 如果返回值不是 ERROR_INSUFFICIENT_BUFFER（122），则返回空字符串
		return ""
	}
	if size <= 0 { // 如果 TCP 连接信息的大小小于等于 0，则返回空字符串
		return ""
	}
	Table := make([]byte, size)                       // 创建一个大小为 size 的字节数组，用于保存 TCP 连接信息
	tcpTablePtr := uintptr(unsafe.Pointer(&Table[0])) // 获取字节数组的指针

	// 调用 syscall.Syscall 函数调用 GetTcpTable2 函数，获取 TCP 连接信息
	ret, _, _ = syscall.Syscall(GetTcpTable2, 3, tcpTablePtr, uintptr(unsafe.Pointer(&size)), 1)

	dwNum := int((*MIB_TCPTABLE2)(unsafe.Pointer(tcpTablePtr)).dwNumEntries) // 获取 TCP 连接信息中记录的数量
	// 遍历 TCP 连接信息中的记录，查找指定的 TCP 连接信息
	for i := 0; i < dwNum; i++ {
		row := (*MIB_TCPROW2)(unsafe.Pointer(tcpTablePtr + 4 + uintptr(i*28))) // 获取当前记录的指针
		// 根据记录中的地址和端口信息，生成本地地址字符串
		LocalAddress := fmt.Sprintf("%d.%d.%d.%d:%d", (row.dwLocalAddr>>0)&0xff, (row.dwLocalAddr>>8)&0xff, (row.dwLocalAddr>>16)&0xff, (row.dwLocalAddr>>24)&0xff, ntohs(row.dwLocalPort))
		if LocalAddress == tcpInfo && ntohs(row.dwRemotePort) == SunnyPort { // 如果找到指定的 TCP 连接信息，则返回 PID
			return strconv.Itoa(int(row.dwOwningPid))
		}
	}
	return "" // 如果没有找到指定的 TCP 连接信息，则返回空字符串
}

// ntohs 用于将 16 位的网络字节序转换为主机字节序
func ntohs(x uint32) int {
	v := uint16(x)
	return int(uint16(byte(v>>8)) | uint16(byte(v))<<8)
}
