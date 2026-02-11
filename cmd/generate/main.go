package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"dove/pkg/generator"
)

func main() {
	var (
		modelFile = flag.String("model", "", "模型文件路径 (例如: internal/model/user.go)")
		outputDir = flag.String("output", "internal", "输出目录")
		help      = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		fmt.Println("CRUD代码生成器")
		fmt.Println("用法: go run cmd/generate/main.go -model <模型文件路径> [-output <输出目录>]")
		fmt.Println("")
		fmt.Println("参数:")
		fmt.Println("  -model    模型文件路径 (必需)")
		fmt.Println("  -output   输出目录 (可选，默认为 internal)")
		fmt.Println("  -help     显示帮助信息")
		return
	}

	if *modelFile == "" {
		log.Fatal("请指定模型文件路径，使用 -model 参数")
	}

	// 检查模型文件是否存在
	if _, err := os.Stat(*modelFile); os.IsNotExist(err) {
		log.Fatalf("模型文件不存在: %s", *modelFile)
	}

	// 解析模型文件
	modelInfo, err := generator.ParseModelFile(*modelFile)
	if err != nil {
		log.Fatalf("解析模型文件失败: %v", err)
	}

	// 获取模型名称（去掉.go扩展名）
	modelName := strings.TrimSuffix(filepath.Base(*modelFile), ".go")
	// 确保模型名称首字母大写
	modelInfo.Name = strings.Title(modelName)

	fmt.Printf("解析模型: %s\n", modelInfo.Name)
	fmt.Printf("字段数量: %d\n", len(modelInfo.Fields))
	fmt.Printf("可搜索字段: %v\n", modelInfo.Searchable)
	fmt.Printf("可排序字段: %v\n", modelInfo.Sortable)

	// 创建生成器
	gen := generator.NewGenerator(modelInfo, *outputDir)

	// 生成所有代码
	if err := gen.GenerateAll(); err != nil {
		log.Fatalf("生成代码失败: %v", err)
	}

	fmt.Println("✅ CRUD代码生成完成!")
	fmt.Printf("输出目录: %s\n", *outputDir)
	fmt.Println("")
	fmt.Println("生成的文件:")
	fmt.Printf("  - internal/domain/%s.go\n", strings.ToLower(modelName))
	fmt.Printf("  - internal/service/%s_service.go\n", strings.ToLower(modelName))
	fmt.Printf("  - internal/repository/%s_repository.go\n", strings.ToLower(modelName))
	fmt.Printf("  - internal/handler/%s_handler.go\n", strings.ToLower(modelName))
	fmt.Println("")
	fmt.Println("注意: 请检查生成的代码并根据需要进行调整")
	fmt.Println("手动补充 provider")
	fmt.Println("手动补充路由")
	fmt.Println("手动 make wire,make swagger")
	fmt.Println("手动补充 migrate")
}
