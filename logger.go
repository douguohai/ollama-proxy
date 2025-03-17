package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogEntry 日志条目结构
type LogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Token       string                 `json:"token,omitempty"`
	TokenValid  bool                   `json:"token_valid,omitempty"`
	RequestBody map[string]interface{} `json:"request_body,omitempty"`
	Response    interface{}            `json:"response,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// Logger 日志记录器
type Logger struct {
	logFile *os.File
}

// NewLogger 创建新的日志记录器
func NewLogger() (*Logger, error) {
	// 确保logs目录存在
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, err
	}

	// 使用当前日期作为日志文件名
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	filePath := filepath.Join(logsDir, fileName)

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{logFile: file}, nil
}

// LogRequest 记录请求和响应
func (l *Logger) LogRequest(method, path string, requestBody map[string]interface{}, response interface{}, err error, token string, tokenValid bool) {
	entry := LogEntry{
		Timestamp:   time.Now().Format(time.RFC3339),
		Method:      method,
		Path:        path,
		Token:       token,
		TokenValid:  tokenValid,
		RequestBody: requestBody,
	}

	if err != nil {
		entry.Error = err.Error()
	} else if response != nil {
		entry.Response = response
	}

	// 将日志条目转换为JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("Error marshaling log entry: %v\n", err)
		return
	}

	// 同时输出到控制台和文件
	fmt.Println(string(jsonData))

	// 写入日志文件
	if _, err := l.logFile.WriteString(string(jsonData) + "\n"); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
