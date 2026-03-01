// Package analyzer 提供日志解析功能，支持从日志行中提取级别、IP等核心信息
// 核心函数为 ParseLog，可解析格式为 "[LEVEL] ... IP: xxx.xxx.xxx.xxx" 的日志行

package analyzer

import (
	"errors"
	"regexp"
)

// 封装日志分析包

// LogEntry 表示一条日志记录
type LogEntry struct {
	Level string // INFO, ERROR, WARN
	IP    string
	Msg   string // 未提取的部分（可后续扩展）
}

// logRegex 是预编译的正则表达式，用于匹配日志格式
var logRegex = regexp.MustCompile(`\[(INFO|ERROR|WARN)\].*?IP:\s*(\d{1,3}(?:\.\d{1,3}){3})`)

// ParseLog 解析日志行，提取日志级别和IP地址
// 返回 LogEntry 指针和错误信息
func ParseLog(line string) (*LogEntry, error) {
	// 使用预编译的正则表达式匹配日志行
	matches := logRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, errors.New("invalid log format")
	}

	// 提取匹配的组：[0]是完整匹配，[1]是级别，[2]是IP
	level := matches[1]
	ip := matches[2]

	// 构建日志条目（Msg 未提取，设为空字符串）
	return &LogEntry{
		Level: level,
		IP:    ip,
		Msg:   "", // 未提取部分留空
	}, nil
}
