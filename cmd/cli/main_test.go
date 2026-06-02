package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCLI_AdminCreate_ShortPassword 测试密码长度不足。
func TestCLI_AdminCreate_ShortPassword(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		createAdminCmd([]string{"-username", "test", "-password", "short"})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCLI_AdminCreate_ShortPassword")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")
	output, _ := cmd.CombinedOutput()
	assert.Contains(t, string(output), "密码长度不足")
}

// TestCLI_AdminCreate_EmptyCredentials 测试用户名和密码为空。
func TestCLI_AdminCreate_EmptyCredentials(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		createAdminCmd([]string{"-username", "", "-password", ""})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCLI_AdminCreate_EmptyCredentials")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")
	output, _ := cmd.CombinedOutput()
	assert.Contains(t, string(output), "用户名和密码不能为空")
}

// TestCLI_UnknownCommand 测试未知命令。
func TestCLI_UnknownCommand(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		os.Args = []string{"mousecake", "unknown"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCLI_UnknownCommand")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")
	output, _ := cmd.CombinedOutput()
	assert.Contains(t, string(output), "Unknown command")
}
