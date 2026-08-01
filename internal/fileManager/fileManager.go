package filemanager

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"text/template"
	"webtest/internal/models"
)

var (
	ctxErr       = errors.New("Отмена запроса со стороны клиента!")
	WriteFileErr = errors.New("Ошибка записи в файл!")
	tmplParseErr = errors.New("Не удалось пропарсить шаблон!")
	tmplExecErr  = errors.New("Ошибка выполнения шаблона!")
)

type ManagerConfig struct {
	RootDir string
	Data    *models.Eat
}

func CreateFile(ctx context.Context, config *ManagerConfig) (string, string, error) {
	fileName := config.Data.Name + ".txt"
	filePath := filepath.Join(config.RootDir, fileName)

	buffer, err := writeData(ctx, config.Data)
	if err != nil {
		return "", "", err
	}

	select {
	case <-ctx.Done():
		return "", "", ctxErr
	default:
		err := os.WriteFile(
			filePath,
			buffer.Bytes(), //записать сюда buffer.Bytes() и проверить с запуском
			0644,
		)
		if err != nil {
			return "", "", WriteFileErr
		}
	}

	return fileName, filePath, nil
}

func writeData(ctx context.Context, data *models.Eat) (*bytes.Buffer, error) {
	select {
	case <-ctx.Done():
		return nil, ctxErr
	default:
		var buffer bytes.Buffer

		text :=
			`
Рецепт "{{.Name}}":

Категория: {{.Category}}

Ингредиенты:
{{range .Ingredients}}- {{.}}
{{end}}
Инструкция:
{{.Instructions}}
`

		tmpl, err := template.New("EatTemplate").Parse(text)
		if err != nil {
			return nil, tmplParseErr
		}

		err = tmpl.Execute(&buffer, data)
		if err != nil {
			return nil, tmplExecErr
		}

		return &buffer, nil
	}
}
