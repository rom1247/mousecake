// Package main 提供 CLI 子命令入口。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/shared/database"
	"github.com/mousecake-go/mousecake-go/internal/user"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mousecake <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "admin":
		adminCmd(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func adminCmd(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: mousecake admin <subcommand>")
		os.Exit(1)
	}

	switch args[0] {
	case "create":
		createAdminCmd(args[1:])
	default:
		fmt.Printf("Unknown admin subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func createAdminCmd(args []string) {
	fs := flag.NewFlagSet("admin create", flag.ExitOnError)
	username := fs.String("username", "", "管理员用户名")
	password := fs.String("password", "", "管理员密码（最少 8 位）")
	fs.Parse(args)

	if *username == "" || *password == "" {
		fmt.Println("用户名和密码不能为空")
		os.Exit(1)
	}

	if len(*password) < 8 {
		fmt.Println("密码长度不足（最少 8 位）")
		os.Exit(1)
	}

	cfg, err := config.Load("config/app.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	repo := user.NewUserRepository(db)
	svc := user.NewService(repo, nil, nil, "", "", "")

	if err := svc.CreateAdmin(context.Background(), *username, *password); err != nil {
		fmt.Printf("创建管理员失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("管理员 %s 创建成功\n", *username)
}
