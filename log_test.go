package xxl

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRotateLog(t *testing.T) {

	// 创建测试用例目录
	testCases := []struct {
		dirName string // 目录名
		age     int    // 目录年龄（天）
		expect  bool   // 期望是否被删除
	}{
		{time.Now().AddDate(0, 0, -31).Format("2006-01-02"), 31, true},  // 31天前，应该被删除
		{time.Now().AddDate(0, 0, -15).Format("2006-01-02"), 15, false}, // 15天前，不应该被删除
		{time.Now().Format("2006-01-02"), 0, false},                     // 今天，不应该被删除
	}

	// 创建测试目录
	for _, tc := range testCases {
		dirPath := filepath.Join(logDir, tc.dirName)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dirPath, err)
		}
		// 在目录中创建一个测试文件
		testFile := filepath.Join(dirPath, "test.log")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFile, err)
		}
	}

	// 创建上下文和取消函数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 启动日志轮转（使用较短的时间间隔进行测试）
	go RotateLog(ctx, 500*time.Millisecond, 30) // 30天的保留期

	// 等待日志轮转执行
	time.Sleep(1 * time.Second)

	// 验证结果
	for _, tc := range testCases {
		dirPath := filepath.Join(logDir, tc.dirName)
		_, err := os.Stat(dirPath)
		exists := !os.IsNotExist(err)

		if exists == tc.expect {
			t.Errorf("Directory %s: expected existence = %v, got %v", tc.dirName, !tc.expect, exists)
		}
	}
}
