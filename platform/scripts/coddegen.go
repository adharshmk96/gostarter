package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Templates for each layer
const (
	handlerTemplate = `package api

import (
    "go.opentelemetry.io/otel/trace"
    "{{.PackageName}}/infra"
    "{{.PackageName}}/internals/domain"
    "log/slog"
    "net/http"
)

type {{.VarName}}Handler struct {
	logger         *slog.Logger
	tracer         trace.Tracer
	{{.VarName}}Service domain.TrainerService
}

func NewTrainerHandler(
	container *infra.Container,
	{{.VarName}}Service domain.TrainerService,
	tokenService domain.TokenService,
) domain.TrainerHandler {
	logger := container.Logger.With("component", "{{.VarName}}Handler")
	return &{{.VarName}}Handler{
		logger:         logger,
		tracer:         container.Tracer,
		{{.VarName}}Service: {{.VarName}}Service,
	}
}
`

	serviceTemplate = `package service

import (
    "{{.PackageName}}/internals/domain"
    "{{.PackageName}}/internals/infra"
    "go.opentelemetry.io/trace"
    "log/slog"
)

type {{.VarName}}Service struct {
    logger *slog.Logger
    tracer trace.Tracer
    {{.VarName}}Repo domain.{{.Name}}Repository
}

func New{{.Name}}Service(
    container *infra.Container,
    {{.VarName}}Repo domain.{{.Name}}Repository,
) domain.{{.Name}}Service {
    return &{{.VarName}}Service{
        logger:      container.Logger,
        tracer:      container.Tracer,
        {{.VarName}}Repo: {{.VarName}}Repo,
    }
}
{{range .Methods}}
func (s *{{$.VarName}}Service) {{.Name}}({{.Params}}) {{.Returns}} {
    // TODO: Implement {{.Name}}
    panic("not implemented")
}
{{end}}`

	repositoryTemplate = `package memory

import (
    "{{.PackageName}}/internals/domain"
    "{{.PackageName}}/internals/infra"
    "go.opentelemetry.io/trace"
    "log/slog"
)

type {{.VarName}}Repository struct {
    logger   *slog.Logger
    tracer   trace.Tracer
    {{.VarName}}s []domain.{{.Name}}
}

func New{{.Name}}Repository(container *infra.Container) domain.{{.Name}}Repository {
    return &{{.VarName}}Repository{
        logger:   container.Logger,
        tracer:   container.Tracer,
        {{.VarName}}s: []domain.{{.Name}}{},
    }
}
{{range .Methods}}
func (r *{{$.VarName}}Repository) {{.Name}}({{.Params}}) {{.Returns}} {
    // TODO: Implement {{.Name}}
    panic("not implemented")
}
{{end}}`
)

type Method struct {
	Name    string
	Params  string
	Returns string
}

type TemplateData struct {
	PackageName string
	Name        string // The entity name (e.g., "Account")
	VarName     string // Lowercase version for variables (e.g., "account")
	Methods     []Method
}

type DomainFile struct {
	Name           string
	HandlerMethods []Method
	ServiceMethods []Method
	RepoMethods    []Method
}

func main() {
	// Find all domain files
	domainFiles, err := scanDomainFiles("internals/domain")
	if err != nil {
		log.Fatalf("Failed to scan domain files: %v", err)
	}

	// Get module name from go.mod
	modName := getModuleName()

	// Generate implementations for each domain file
	for _, df := range domainFiles {
		generateImplementations(modName, df)
	}
}

func scanDomainFiles(domainDir string) ([]DomainFile, error) {
	var domainFiles []DomainFile

	err := filepath.Walk(domainDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// Parse the file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %v", path, err)
		}

		// Extract domain name from filename
		baseName := strings.TrimSuffix(info.Name(), ".go")
		entityName := strings.Title(baseName) // Convert first letter to uppercase

		var df DomainFile
		df.Name = entityName

		// Find interfaces
		ast.Inspect(node, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}

			switch {
			case strings.HasSuffix(typeSpec.Name.Name, "Handler"):
				df.HandlerMethods = extractMethods(interfaceType)
			case strings.HasSuffix(typeSpec.Name.Name, "Service"):
				df.ServiceMethods = extractMethods(interfaceType)
			case strings.HasSuffix(typeSpec.Name.Name, "Repository"):
				df.RepoMethods = extractMethods(interfaceType)
			}

			return true
		})

		if len(df.HandlerMethods) > 0 || len(df.ServiceMethods) > 0 || len(df.RepoMethods) > 0 {
			domainFiles = append(domainFiles, df)
		}

		return nil
	})

	return domainFiles, err
}

func generateImplementations(modName string, df DomainFile) {
	data := TemplateData{
		PackageName: modName,
		Name:        df.Name,
		VarName:     strings.ToLower(df.Name),
	}

	// Generate handler
	if len(df.HandlerMethods) > 0 {
		data.Methods = df.HandlerMethods
		generateFile(
			fmt.Sprintf("internals/delivery/http/api/%s.go", strings.ToLower(df.Name)),
			handlerTemplate,
			data,
		)
	}

	// Generate service
	if len(df.ServiceMethods) > 0 {
		data.Methods = df.ServiceMethods
		generateFile(
			fmt.Sprintf("internals/service/%s.go", strings.ToLower(df.Name)),
			serviceTemplate,
			data,
		)
	}

	// Generate repository
	if len(df.RepoMethods) > 0 {
		data.Methods = df.RepoMethods
		generateFile(
			fmt.Sprintf("internals/storage/memory/%s.go", strings.ToLower(df.Name)),
			repositoryTemplate,
			data,
		)
	}
}

func extractMethods(interfaceType *ast.InterfaceType) []Method {
	var methods []Method
	for _, method := range interfaceType.Methods.List {
		funcType := method.Type.(*ast.FuncType)

		params := extractParams(funcType.Params)
		returns := extractReturns(funcType.Results)

		methods = append(methods, Method{
			Name:    method.Names[0].Name,
			Params:  params,
			Returns: returns,
		})
	}
	return methods
}

func extractParams(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	var params []string
	for _, field := range fields.List {
		paramType := ""
		switch t := field.Type.(type) {
		case *ast.Ident:
			paramType = t.Name
		case *ast.SelectorExpr:
			paramType = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
		case *ast.StarExpr:
			switch x := t.X.(type) {
			case *ast.Ident:
				paramType = "*" + x.Name
			case *ast.SelectorExpr:
				paramType = fmt.Sprintf("*%s.%s", x.X.(*ast.Ident).Name, x.Sel.Name)
			}
		}

		if len(field.Names) == 0 {
			params = append(params, paramType)
		} else {
			for _, name := range field.Names {
				params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
			}
		}
	}
	return strings.Join(params, ", ")
}

func extractReturns(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	var returns []string
	for _, field := range fields.List {
		returnType := ""
		switch t := field.Type.(type) {
		case *ast.Ident:
			returnType = t.Name
		case *ast.StarExpr:
			switch x := t.X.(type) {
			case *ast.Ident:
				returnType = "*" + x.Name
			case *ast.SelectorExpr:
				returnType = fmt.Sprintf("*%s.%s", x.X.(*ast.Ident).Name, x.Sel.Name)
			}
		case *ast.SelectorExpr:
			returnType = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
		}
		returns = append(returns, returnType)
	}

	if len(returns) == 0 {
		return ""
	}
	if len(returns) == 1 {
		return returns[0]
	}
	return "(" + strings.Join(returns, ", ") + ")"
}

func generateFile(path string, tmpl string, data TemplateData) {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		log.Printf("File %s already exists, skipping", path)
		return
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Parse and execute template
	t, err := template.New("file").Parse(tmpl)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	// Write file
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		log.Fatalf("Failed to write file %s: %v", path, err)
	}

	fmt.Printf("Generated %s\n", path)
}

func getModuleName() string {
	// Read go.mod file
	data, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatal("Failed to read go.mod file")
	}

	// Extract module name from first line
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		log.Fatal("Empty go.mod file")
	}

	parts := strings.Fields(lines[0])
	if len(parts) != 2 || parts[0] != "module" {
		log.Fatal("Invalid go.mod file format")
	}

	return parts[1]
}
