package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ModelInfo 模型信息
type ModelInfo struct {
	Name       string   // 模型名称
	Package    string   // 包名
	Fields     []Field  // 字段列表
	TableName  string   // 表名
	PrimaryKey string   // 主键字段
	Searchable []string // 可搜索字段
	Sortable   []string // 可排序字段
}

// Field 字段信息
type Field struct {
	Name     string // 字段名
	Type     string // 字段类型
	Tag      string // 标签
	JSONName string // JSON 名称
	Comment  string // 注释
}

// Generator CRUD 代码生成器
type Generator struct {
	ModelInfo *ModelInfo
	OutputDir string
}

// NewGenerator 创建生成器
func NewGenerator(modelInfo *ModelInfo, outputDir string) *Generator {
	return &Generator{
		ModelInfo: modelInfo,
		OutputDir: outputDir,
	}
}

// GenerateAll 生成所有 CRUD 代码
func (g *Generator) GenerateAll() error {
	// 创建输出目录
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// 生成 domain 层代码
	if err := g.generateDomain(); err != nil {
		return fmt.Errorf("failed to generate domain: %v", err)
	}

	// 生成 service 层代码
	if err := g.generateService(); err != nil {
		return fmt.Errorf("failed to generate service: %v", err)
	}

	// 生成 repository 层代码
	if err := g.generateRepository(); err != nil {
		return fmt.Errorf("failed to generate repository: %v", err)
	}

	// 生成 handler 层代码
	if err := g.generateHandler(); err != nil {
		return fmt.Errorf("failed to generate handler: %v", err)
	}

	// 更新 wire 配置
	if err := g.updateWire(); err != nil {
		return fmt.Errorf("failed to update wire: %v", err)
	}

	return nil
}

// generateDomain 生成 domain 层代码
func (g *Generator) generateDomain() error {
	// 创建函数映射
	funcMap := template.FuncMap{
		"lower":    strings.ToLower,
		"title":    strings.ToTitle,
		"contains": strings.Contains,
	}

	// 创建模板并设置函数映射
	tmpl := template.New("domain.go.tmpl").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("pkg/generator/templates/domain.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse domain template: %v", err)
	}

	outputFile := filepath.Join(g.OutputDir, "domain", strings.ToLower(g.ModelInfo.Name)+".go")
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create domain directory: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create domain file: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, g.ModelInfo)
}

// generateService 生成 service 层代码
func (g *Generator) generateService() error {
	// 创建函数映射
	funcMap := template.FuncMap{
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"contains": strings.Contains,
	}

	// 创建模板并设置函数映射
	tmpl := template.New("service.go.tmpl").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("pkg/generator/templates/service.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse service template: %v", err)
	}

	outputFile := filepath.Join(g.OutputDir, "service", strings.ToLower(g.ModelInfo.Name)+"_service.go")
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, g.ModelInfo)
}

// generateRepository 生成 repository 层代码
func (g *Generator) generateRepository() error {
	// 创建函数映射
	funcMap := template.FuncMap{
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"contains": strings.Contains,
	}

	// 创建模板并设置函数映射
	tmpl := template.New("repository.go.tmpl").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("pkg/generator/templates/repository.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %v", err)
	}

	outputFile := filepath.Join(g.OutputDir, "repository", strings.ToLower(g.ModelInfo.Name)+"_repository.go")
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create repository file: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, g.ModelInfo)
}

// generateHandler 生成 handler 层代码
func (g *Generator) generateHandler() error {
	// 创建函数映射
	funcMap := template.FuncMap{
		"lower":    strings.ToLower,
		"title":    strings.Title,
		"contains": strings.Contains,
	}

	// 创建模板并设置函数映射
	tmpl := template.New("handler.go.tmpl").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("pkg/generator/templates/handler.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse handler template: %v", err)
	}

	outputFile := filepath.Join(g.OutputDir, "handler", strings.ToLower(g.ModelInfo.Name)+"_handler.go")
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create handler directory: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create handler file: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, g.ModelInfo)
}

// updateWire 更新 wire 配置
func (g *Generator) updateWire() error {
	wireFile := filepath.Join(g.OutputDir, "wire", "providers.go")
	if _, err := os.Stat(wireFile); os.IsNotExist(err) {
		// 如果wire文件不存在，创建一个新的
		return g.createWireFile(wireFile)
	}

	// 读取现有的wire文件
	content, err := os.ReadFile(wireFile)
	if err != nil {
		return fmt.Errorf("failed to read wire file: %v", err)
	}

	// 检查是否已经包含了新的provider
	if strings.Contains(string(content), "New"+g.ModelInfo.Name+"Repository") {
		// 已经存在，不需要更新
		return nil
	}

	// 更新wire文件
	return g.updateWireFile(wireFile, string(content))
}

// createWireFile 创建新的wire文件
func (g *Generator) createWireFile(wireFile string) error {
	tmpl := `package wire

import (
	"dove/internal/app"
	"dove/internal/handler"
	"dove/internal/repository"
	"dove/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet 是 wire 的提供者集合
var ProviderSet = wire.NewSet(
	// Repository 层
	repository.New{{.Name}}Repository,

	// Service 层
	service.New{{.Name}}Service,

	// Handler 层
	handler.New{{.Name}}Handler,

	// 提供 gin 引擎
	ProvideGinEngine,

	// 提供应用实例
	app.NewApp,
	app.InitializeDatabase,
)

// ProvideGinEngine 提供 gin 引擎
func ProvideGinEngine() *gin.Engine {
	engine := gin.New()
	return engine
}
`

	template, err := template.New("wire").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse wire template: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(wireFile), 0755); err != nil {
		return fmt.Errorf("failed to create wire directory: %v", err)
	}

	file, err := os.Create(wireFile)
	if err != nil {
		return fmt.Errorf("failed to create wire file: %v", err)
	}
	defer file.Close()

	return template.Execute(file, g.ModelInfo)
}

// updateWireFile 更新现有的wire文件
func (g *Generator) updateWireFile(wireFile, content string) error {
	// 简单的字符串替换来添加新的provider
	modelName := g.ModelInfo.Name

	// 在repository部分添加新的repository
	repoPattern := `	repository.NewProductRepository,
	repository.NewUserRepository,`
	repoReplacement := `	repository.NewProductRepository,
	repository.NewUserRepository,
	repository.New` + modelName + `Repository,`

	// 在service部分添加新的service
	servicePattern := `	service.NewProductService,
	service.NewUserService,`
	serviceReplacement := `	service.NewProductService,
	service.NewUserService,
	service.New` + modelName + `Service,`

	// 在handler部分添加新的handler
	handlerPattern := `	handler.NewProductHandler,
	handler.NewUserHandler,
	handler.NewAuthHandler,
	handler.NewHealthHandler,`
	handlerReplacement := `	handler.NewProductHandler,
	handler.NewUserHandler,
	handler.New` + modelName + `Handler,
	handler.NewAuthHandler,
	handler.NewHealthHandler,`

	// 执行替换
	content = strings.Replace(content, repoPattern, repoReplacement, 1)
	content = strings.Replace(content, servicePattern, serviceReplacement, 1)
	content = strings.Replace(content, handlerPattern, handlerReplacement, 1)

	// 写回文件
	return os.WriteFile(wireFile, []byte(content), 0644)
}

// ParseModelFile 解析模型文件
func ParseModelFile(filePath string) (*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	var modelInfo ModelInfo
	modelInfo.Package = node.Name.Name

	// 遍历 AST 查找结构体定义
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						modelInfo.Name = typeSpec.Name.Name
						modelInfo.TableName = strings.ToLower(modelInfo.Name) + "s"
						modelInfo.PrimaryKey = "id"

						// 解析字段
						for _, field := range structType.Fields.List {
							if len(field.Names) > 0 {
								fieldName := field.Names[0].Name
								fieldType := getTypeString(field.Type)
								fieldTag := ""
								if field.Tag != nil {
									fieldTag = strings.Trim(field.Tag.Value, "`")
								}

								// 解析 JSON 标签
								jsonName := fieldName
								if fieldTag != "" {
									if strings.Contains(fieldTag, `json:"`) {
										jsonStart := strings.Index(fieldTag, `json:"`) + 6
										jsonEnd := strings.Index(fieldTag[jsonStart:], `"`)
										if jsonEnd > 0 {
											jsonName = fieldTag[jsonStart : jsonStart+jsonEnd]
										}
									}
								}

								// 获取注释
								comment := ""
								if field.Doc != nil {
									comment = strings.TrimSpace(field.Doc.Text())
								}

								fieldInfo := Field{
									Name:     fieldName,
									Type:     fieldType,
									Tag:      fieldTag,
									JSONName: jsonName,
									Comment:  comment,
								}

								modelInfo.Fields = append(modelInfo.Fields, fieldInfo)

								// 判断是否为可搜索字段
								if isSearchableField(fieldName) {
									modelInfo.Searchable = append(modelInfo.Searchable, fieldName)
								}

								// 判断是否为可排序字段
								if isSortableField(fieldName) {
									modelInfo.Sortable = append(modelInfo.Sortable, fieldName)
								}
							}
						}
					}
				}
			}
		}
	}

	return &modelInfo, nil
}

// getTypeString 获取类型字符串
func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.ArrayType:
		return "[]" + getTypeString(t.Elt)
	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name
	default:
		return "interface{}"
	}
}

// isSearchableField 判断是否为可搜索字段
func isSearchableField(fieldName string) bool {
	searchableFields := []string{"name", "title", "description", "username", "email", "nickname"}
	for _, field := range searchableFields {
		if strings.Contains(strings.ToLower(fieldName), field) {
			return true
		}
	}
	return false
}

// isSortableField 判断是否为可排序字段
func isSortableField(fieldName string) bool {
	sortableFields := []string{"id", "name", "created_at", "updated_at", "price", "stock", "status"}
	for _, field := range sortableFields {
		if strings.Contains(strings.ToLower(fieldName), field) {
			return true
		}
	}
	return false
}
