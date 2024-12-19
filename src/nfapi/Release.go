//go:build windows
// +build windows

package NFapi

/*
#include <windows.h>

BOOL disableWow64FsRedirection(PVOID* oldValue) {
    return Wow64DisableWow64FsRedirection(oldValue);
}

BOOL revertWow64FsRedirection(PVOID oldValue) {
    return Wow64RevertWow64FsRedirection(oldValue);
}
*/
import "C"
import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/Resource"
	"github.com/qtgolang/SunnyNet/src/public"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// 生成指定长度的随机字母串
func randomLetters(length int) string {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成指定长度的随机字母
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// 删除旧的驱动文件
func deleteOldFiles() {
	OldFileName := System32Dir + "\\drivers\\SunnyFilter.sys"
	//复制到临时目录去系统重启后才可删除
	_ = MoveFileToTempDir(OldFileName, "Sunny_"+randomLetters(32)+extensionsTemp)
	//删除临时目录下的所有sys 文件
	tempDir := os.TempDir()
	// 搜索所有 .sys 文件
	_ = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查文件是否是 .sys 文件
		if !info.IsDir() && filepath.Ext(path) == extensionsTemp {
			_ = os.Remove(path)
		}
		return nil
	})
}

// MoveFileToTempDir 将指定文件移动到 Windows 临时目录
// srcFile: 源文件路径，destFileName：目标文件名
// 返回值：目标文件的路径，以及可能出现的错误
func MoveFileToTempDir(srcFile, destFileName string) string {
	tempDir := os.TempDir()
	// 拼接目标文件路径
	destPath := filepath.Join(tempDir, destFileName)
	// 移动文件
	err := os.Rename(srcFile, destPath)
	if err != nil {
		return ""
	}
	return destPath
}

// System32Dir C:\Windows\system32\
var System32Dir = GetSystemDirectory()
var extensionsTemp = ".tmpSys"
var DriverFile = System32Dir + "\\drivers\\" + NF_DriverName + ".sys"

func Install() string {
	deleteOldFiles()
	//XP直接打开不程序，所以就直接忽略
	s := []string{"OS", "Get", "Caption"}
	IsWin7 := strings.Index(ExecCommand("Wmic", s), "Windows 7") != -1
	var oldValue uintptr
	if Is64Windows {
		//如果是32位进程 禁止文件重定向 驱动只能写到 system32 目录
		if !WindowsX64 {
			oldValue = Wow64DisableWow64FsRedirection()
		}
	}
	if !Exists(DriverFile) {
		if IsWin7 {
			if Is64Windows {
				WriteFile(DriverFile, Resource.TdiAmd64Netfilter2)

			} else {
				WriteFile(DriverFile, Resource.TdiI386Netfilter2)
			}
		} else {
			if Is64Windows {
				WriteFile(DriverFile, Resource.WfpAmd64Netfilter2)
			} else {
				WriteFile(DriverFile, Resource.WfpI386Netfilter2)
			}
		}
	}
	if Is64Windows {
		//如果是32位进程 恢复文件重定向
		if !WindowsX64 {
			Wow64RevertWow64FsRedirection(oldValue)
		}

		WriteFile("DrDLL", Resource.NfapiX64Nfapi)
	}

	DrDLL := ""
	if WindowsX64 {
		DrDLL = WindowsDirectory + NF_DLLName + "64.dll"
		WriteFile(DrDLL, Resource.NfapiX64Nfapi)
	} else {
		DrDLL = WindowsDirectory + NF_DLLName + "32.dll"
		WriteFile(DrDLL, Resource.NfapiWin32Nfapi)
	}
	return DrDLL
}

var (
	WindowsDirectory = GetWindowsDirectory()
	NspPath          = WindowsDirectory + "SunnyNet_Nsp.dll"
)

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true

}

// Wow64DisableWow64FsRedirection 禁用调用线程的文件系统重定向，默认情况下启用文件系统重定向。此功能对于想要访问本机system32目录的32位应用程序很有用。
func Wow64DisableWow64FsRedirection() uintptr {
	var oldValue C.PVOID
	success := C.disableWow64FsRedirection(&oldValue)
	if success == 0 {
		fmt.Println("禁用文件系统重定向 失败")
		return 0
	}
	return uintptr(oldValue)
}

// Wow64RevertWow64FsRedirection 恢复调用线程的文件系统重定向。
func Wow64RevertWow64FsRedirection(oldValue uintptr) bool {
	success := 0
	if oldValue == 0 {
		var oldValues C.PVOID
		success = int(C.revertWow64FsRedirection(oldValues))
	} else {
		success = int(C.revertWow64FsRedirection(C.PVOID(oldValue)))
	}
	if success == 0 {
		fmt.Println("恢复文件系统重定向 失败")
		return false
	}

	return true
}
func WriteFile(path string, data []byte) {
	if checkFileIsExist(path) {
		err := os.Remove(path)
		if err != nil {
			return
		}
	}
	f, err1 := os.Create(path) //创建文件
	if err1 == nil {
		_, err1 = f.Write(data)
		if err1 != nil {

			return
		}
		err1 = f.Close()
		if err1 != nil {

			return
		}
	} else {
		if err1 != nil {
			return
		}
	}
}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func ExecCommand(commandName string, params []string) string {
	cmd := exec.Command(commandName, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	_ = cmd.Start()
	var s []byte
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadBytes('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		s = public.BytesCombine(s, line)
	}
	return string(s)
}
