package Api

import "github.com/qtgolang/SunnyNet/SunnyNet"

const (
	//OtherCommandDisable_TCP 禁用TCP 返回0失败 返回1成功
	OtherCommandDisable_TCP = uintptr(1001)
)

func OtherCommands(Cmd uintptr, Command ...uintptr) uintptr {
	switch Cmd {
	case OtherCommandDisable_TCP:
		{
			if len(Command) < 1 {
				return 0
			}
			SunnyContext := int(Command[0])
			disable := int(Command[1]) == 1
			SunnyNet.SunnyStorageLock.Lock()
			w := SunnyNet.SunnyStorage[SunnyContext]
			SunnyNet.SunnyStorageLock.Unlock()
			if w == nil {
				return 0
			}
			w.DisableTCP(disable)
		}
		return 1

	}
	return 0
}
