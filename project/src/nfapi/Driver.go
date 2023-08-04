package NFapi

/*
#include "Driver.h"
*/
import "C"
import (
	"SunnyNet/project/public"
	"unsafe"
)

//export go_threadStart
func go_threadStart() {
	threadStart()
}

//export go_threadEnd
func go_threadEnd() {
	threadEnd()
}

//export go_tcpConnectRequest
func go_tcpConnectRequest(id C.ulonglong, pConnInfo uintptr) {
	if pConnInfo == 0 {
		return
	}
	A := (*NF_TCP_CONN_INFO)(unsafe.Pointer(pConnInfo))
	tcpConnectRequest(uint64(id), A)
}

//export go_tcpConnected
func go_tcpConnected(id C.ulonglong, pConnInfo uintptr) {
	if pConnInfo == 0 {
		return
	}
	A := (*NF_TCP_CONN_INFO)(unsafe.Pointer(pConnInfo))
	tcpConnected(uint64(id), A)
}

//export go_tcpClosed
func go_tcpClosed(id C.ulonglong, pConnInfo uintptr) {
	if pConnInfo == 0 {
		return
	}
	A := (*NF_TCP_CONN_INFO)(unsafe.Pointer(pConnInfo))
	tcpClosed(uint64(id), A)
}

//export go_tcpReceive
func go_tcpReceive(id C.ulonglong, buf *byte, len C.int) {
	tcpReceive(uint64(id), buf, int32(len))
}

//export go_tcpSend
func go_tcpSend(id C.ulonglong, buf *byte, len C.int) {
	tcpSend(uint64(id), buf, int32(len))
}

//export go_tcpCanReceive
func go_tcpCanReceive(id C.ulonglong) {
	tcpCanReceive(uint64(id))
}

//export go_tcpCanSend
func go_tcpCanSend(id C.ulonglong) {
	tcpCanSend(uint64(id))
}

//export go_udpCreated
func go_udpCreated(id C.ulonglong, pConnInfo uintptr) {
	if pConnInfo == 0 {
		return
	}
	A := (*NF_UDP_CONN_INFO)(unsafe.Pointer(pConnInfo))
	udpCreated(uint64(id), A)
}

//export go_udpConnectRequest
func go_udpConnectRequest(id C.ulonglong, pConnReq uintptr) {
	if pConnReq == 0 {
		return
	}
	A := (*NF_UDP_CONN_REQUEST)(unsafe.Pointer(pConnReq))
	udpConnectRequest(uint64(id), A)
}

//export go_udpClosed
func go_udpClosed(id C.ulonglong, pConnInfo uintptr) {
	if pConnInfo == 0 {
		return
	}
	A := (*NF_UDP_CONN_INFO)(unsafe.Pointer(pConnInfo))
	udpClosed(uint64(id), A)
}

//export go_udpReceive
func go_udpReceive(id C.ENDPOINT_ID, remoteAddress uintptr, buf uintptr, length C.int, options uintptr) {
	bs := public.CStringToBytes(buf, int(length))
	if remoteAddress == 0 || options == 0 {
		return
	}
	A := (*SockaddrInx)(unsafe.Pointer(remoteAddress))
	B := (*NF_UDP_OPTIONS)(unsafe.Pointer(options))
	udpReceive(uint64(id), A, bs, B)
}

//export go_udpSend
func go_udpSend(id C.ENDPOINT_ID, remoteAddress uintptr, buf uintptr, length C.int, options uintptr) {
	bs := public.CStringToBytes(buf, int(length))
	if remoteAddress == 0 || options == 0 {
		return
	}
	A := (*SockaddrInx)(unsafe.Pointer(remoteAddress))
	B := (*NF_UDP_OPTIONS)(unsafe.Pointer(options))
	udpSend(uint64(id), A, bs, B)
}

//export go_udpCanReceive
func go_udpCanReceive(id C.ulonglong) {
	udpCanReceive(uint64(id))
}

//export go_udpCanSend
func go_udpCanSend(id C.ulonglong) {
	udpCanSend(uint64(id))
}

//******************************************************

func CgoDriverInit(driverName string, InitAddr uintptr) C.int {
	return C.NfDriverInit(C.CString(driverName), unsafe.Pointer(InitAddr))
}
