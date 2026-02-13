package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/deantook/dove/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	configPath    = flag.String("config", "configs/config.yaml", "配置文件路径")
	migrationsDir = flag.String("path", "migrations", "迁移文件目录")
	command       = flag.String("command", "", "迁移命令: up, down, force, version, create")
	version       = flag.Int("version", 0, "版本号 (用于 down 和 force 命令)")
	name          = flag.String("name", "", "迁移文件名 (用于 create 命令)")
)

func main() {
	flag.Parse()

	if *command == "" {
		printUsage()
		os.Exit(1)
	}

	switch *command {
	case "up":
		runUp()
	case "down":
		runDown()
	case "force":
		runForce()
	case "version":
		runVersion()
	case "create":
		runCreate()
	default:
		fmt.Printf("未知命令: %s\n", *command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("数据库迁移工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  migrate -command=<command> [选项]")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  up        - 执行所有待执行的迁移")
	fmt.Println("  down      - 回滚一次迁移")
	fmt.Println("  down N    - 回滚 N 次迁移 (使用 -version=N)")
	fmt.Println("  force     - 强制设置迁移版本 (使用 -version=N)")
	fmt.Println("  version   - 显示当前迁移版本")
	fmt.Println("  create    - 创建新的迁移文件 (使用 -name=<name>)")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -config   - 配置文件路径 (默认: configs/config.yaml)")
	fmt.Println("  -path     - 迁移文件目录 (默认: migrations)")
	fmt.Println("  -version  - 版本号 (用于 down 和 force 命令)")
	fmt.Println("  -name     - 迁移文件名 (用于 create 命令)")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  migrate -command=up")
	fmt.Println("  migrate -command=down")
	fmt.Println("  migrate -command=down -version=1")
	fmt.Println("  migrate -command=force -version=1")
	fmt.Println("  migrate -command=version")
	fmt.Println("  migrate -command=create -name=add_user_table")
}

func getMigrate() (*migrate.Migrate, error) {
	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 获取迁移文件目录的绝对路径
	migrationsPath, err := filepath.Abs(*migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("获取迁移目录路径失败: %w", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("迁移目录不存在: %s", migrationsPath)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// 构建文件系统路径
	fileURL := fmt.Sprintf("file://%s", migrationsPath)

	// 创建 migrate 实例
	m, err := migrate.New(fileURL, dsn)
	if err != nil {
		return nil, fmt.Errorf("创建迁移实例失败: %w", err)
	}

	return m, nil
}

func runUp() {
	m, err := getMigrate()
	if err != nil {
		log.Fatalf("初始化迁移失败: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("没有待执行的迁移")
			return
		}
		log.Fatalf("执行迁移失败: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("获取版本失败: %v", err)
	}

	fmt.Printf("迁移成功完成，当前版本: %d (dirty: %v)\n", version, dirty)
}

func runDown() {
	m, err := getMigrate()
	if err != nil {
		log.Fatalf("初始化迁移失败: %v", err)
	}
	defer m.Close()

	if *version > 0 {
		err = m.Migrate(uint(*version))
	} else {
		err = m.Down()
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("没有可回滚的迁移")
			return
		}
		log.Fatalf("回滚迁移失败: %v", err)
	}

	ver, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("迁移已回滚到初始状态")
			return
		}
		log.Fatalf("获取版本失败: %v", err)
	}

	fmt.Printf("回滚成功，当前版本: %d (dirty: %v)\n", ver, dirty)
}

func runForce() {
	if *version < 0 {
		log.Fatal("版本号必须 >= 0")
	}

	m, err := getMigrate()
	if err != nil {
		log.Fatalf("初始化迁移失败: %v", err)
	}
	defer m.Close()

	if err := m.Force(int(*version)); err != nil {
		log.Fatalf("强制设置版本失败: %v", err)
	}

	fmt.Printf("版本已强制设置为: %d\n", *version)
}

func runVersion() {
	m, err := getMigrate()
	if err != nil {
		log.Fatalf("初始化迁移失败: %v", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("当前版本: 无 (数据库未初始化)")
			return
		}
		log.Fatalf("获取版本失败: %v", err)
	}

	fmt.Printf("当前版本: %d (dirty: %v)\n", version, dirty)
}

func runCreate() {
	if *name == "" {
		log.Fatal("请使用 -name 指定迁移文件名")
	}

	// 获取迁移目录的绝对路径
	migrationsPath, err := filepath.Abs(*migrationsDir)
	if err != nil {
		log.Fatalf("获取迁移目录路径失败: %v", err)
	}

	// 创建目录（如果不存在）
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		log.Fatalf("创建迁移目录失败: %v", err)
	}

	// 获取下一个版本号
	version := getNextVersion(migrationsPath)

	// 格式化版本号（6位数字，前面补0）
	versionStr := fmt.Sprintf("%06d", version)

	// 创建 up 文件
	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.up.sql", versionStr, *name))
	if err := createMigrationFile(upFile, "-- Up Migration\n"); err != nil {
		log.Fatalf("创建 up 文件失败: %v", err)
	}
	fmt.Printf("创建文件: %s\n", upFile)

	// 创建 down 文件
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.down.sql", versionStr, *name))
	if err := createMigrationFile(downFile, "-- Down Migration\n"); err != nil {
		log.Fatalf("创建 down 文件失败: %v", err)
	}
	fmt.Printf("创建文件: %s\n", downFile)

	fmt.Printf("迁移文件创建成功，版本号: %s\n", versionStr)
}

func createMigrationFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func getNextVersion(migrationsPath string) int {
	// 读取目录中的所有文件
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		// 如果目录不存在或读取失败，返回版本 1
		return 1
	}

	maxVersion := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		// 检查是否是迁移文件格式: {version}_{name}.up.sql 或 {version}_{name}.down.sql
		if !(strings.HasSuffix(name, ".up.sql") || strings.HasSuffix(name, ".down.sql")) {
			continue
		}

		// 提取版本号（文件名第一部分）
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}

		// 解析版本号
		var version int
		if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
			continue
		}

		if version > maxVersion {
			maxVersion = version
		}
	}

	return maxVersion + 1
}
