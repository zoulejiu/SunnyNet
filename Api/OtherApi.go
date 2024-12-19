package Api

import (
	"bytes"
	"encoding/binary"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"github.com/qtgolang/SunnyNet/src/public"
	"strings"
)
 
func HexDump(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	hexStr := strings.ReplaceAll(hex.Dump(bin), "\n", "\r\n")
	return public.PointerPtr(hexStr)
}

// BytesToInt 将Go int的Bytes 转为int
func BytesToInt(data uintptr, dataLen int) int {
	bys := public.CStringToBytes(data, dataLen)
	buff := bytes.NewBuffer(bys)
	var B int64
	_ = binary.Read(buff, binary.BigEndian, &B)
	return int(B)
}
