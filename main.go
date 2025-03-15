package main

import "C"
import (
	"bufio"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/http"
	_ "github.com/qtgolang/SunnyNet/src/http/pprof"
	"os"
	"runtime"
	"strings"
)

func init() {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:6001", nil)
	}()
}

func main() {
	Test()
}

// HostEntry 表示一个 hosts 文件中的条目
type HostEntry struct {
	IP        string   // IP 地址
	Hostnames []string // 域名列表
	RawLine   string   // 原始行内容（如果是注释或空行）
}

// ReadAndParseHosts 读取并解析 hosts 文件
func ReadAndParseHosts() ([]HostEntry, error) {
	filePath := getHostsFilePath()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	var entries []HostEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 忽略空行或注释
		if line == "" || strings.HasPrefix(line, "#") {
			entries = append(entries, HostEntry{RawLine: line})
			continue
		}

		// 拆分 IP 与域名
		fields := strings.Fields(line)
		if len(fields) < 2 {
			entries = append(entries, HostEntry{RawLine: line})
			continue
		}

		ip := fields[0]
		hostnames := fields[1:]

		entry := HostEntry{
			IP:        ip,
			Hostnames: hostnames,
			RawLine:   line,
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件出错: %v", err)
	}

	return entries, nil
}

// getHostsFilePath 返回操作系统对应的 hosts 文件路径
func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return `/etc/hosts`
}

func mainxx() {
	entries, err := ReadAndParseHosts()
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	// 打印读取到的所有条目
	for _, entry := range entries {
		if entry.IP == "" {
			//fmt.Println(entry.RawLine) // 打印注释或空行
		} else {
			fmt.Printf("IP: %s, Hostnames: %v\n", entry.IP, entry.Hostnames)
		}
	}
}
