package main

import (
	"flag"
	"log"
	"os"

	"dove/pkg/config"
	"dove/pkg/database"
	"dove/pkg/migrate"
)

func main() {
	// 解析命令行参数
	action := flag.String("action", "migrate", "Action to perform: migrate, reset, drop")
	env := flag.String("env", "dev", "Environment: dev, test, production")
	flag.Parse()

	// 加载配置文件
	if err := config.LoadConfig(*env); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 根据参数执行相应操作
	switch *action {
	case "migrate":
		log.Println("Running migrations...")
		if err := migrate.RunMigrations(); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Migrations completed successfully")

	case "reset":
		log.Println("Resetting database...")
		if err := migrate.ResetDatabase(); err != nil {
			log.Fatal("Failed to reset database:", err)
		}
		log.Println("Database reset completed successfully")

	case "drop":
		log.Println("Dropping all tables...")
		if err := migrate.DropTables(); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		log.Println("All tables dropped successfully")

	default:
		log.Printf("Unknown action: %s", *action)
		log.Println("Available actions: migrate, reset, drop")
		os.Exit(1)
	}
}
