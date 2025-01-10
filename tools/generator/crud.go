// tools/generator/crud.go
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Field struct {
	Name    string // ชื่อฟิลด์
	Type    string // ประเภทข้อมูล
	JsonTag string // json tag
}

type Entity struct {
	Name   string  // ชื่อ Entity
	Fields []Field // ฟิลด์ทั้งหมด
}

// Templates
var domainTmpl = `package domain

type {{.Name}} struct {
    ID        string    ` + "`json:\"id\"`" + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    {{range .Fields}}
    {{.Name}} {{.Type}} ` + "`json:\"{{.JsonTag}}\"`" + `
    {{end}}
}
`

var repoTmpl = `package repository

type {{.Name}}Repository interface {
    Create(ctx context.Context, {{.Name | lower}} *domain.{{.Name}}) error
    GetByID(ctx context.Context, id string) (*domain.{{.Name}}, error)
    List(ctx context.Context) ([]*domain.{{.Name}}, error)
    Update(ctx context.Context, {{.Name | lower}} *domain.{{.Name}}) error
    Delete(ctx context.Context, id string) error
}
`

var usecaseTmpl = `package usecase

type {{.Name}}UseCase struct {
    repo repository.{{.Name}}Repository
}

func New{{.Name}}UseCase(repo repository.{{.Name}}Repository) *{{.Name}}UseCase {
    return &{{.Name}}UseCase{repo: repo}
}

func (u *{{.Name}}UseCase) Create(ctx context.Context, input *domain.{{.Name}}) error {
    return u.repo.Create(ctx, input)
}
`

var handlerTmpl = `package handler

type {{.Name}}Handler struct {
    useCase *usecase.{{.Name}}UseCase
}

func New{{.Name}}Handler(useCase *usecase.{{.Name}}UseCase) *{{.Name}}Handler {
    return &{{.Name}}Handler{useCase: useCase}
}

func (h *{{.Name}}Handler) Create(c *fiber.Ctx) error {
    var input domain.{{.Name}}
    if err := c.BodyParser(&input); err != nil {
        return err
    }
    return h.useCase.Create(c.Context(), &input)
}
`

func Generate(e Entity) error {
	// สร้าง directories
	dirs := []string{
		"internal/domain",
		"internal/repository",
		"internal/usecase",
		"internal/handler",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Generate files
	files := map[string]string{
		filepath.Join("internal/domain", fmt.Sprintf("%s.go", e.Name)):     domainTmpl,
		filepath.Join("internal/repository", fmt.Sprintf("%s.go", e.Name)): repoTmpl,
		filepath.Join("internal/usecase", fmt.Sprintf("%s.go", e.Name)):    usecaseTmpl,
		filepath.Join("internal/handler", fmt.Sprintf("%s.go", e.Name)):    handlerTmpl,
	}

	for path, tmpl := range files {
		if err := generateFile(path, tmpl, e); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(path, tmpl string, data Entity) error {
	t, err := template.New("entity").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}
