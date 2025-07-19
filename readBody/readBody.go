package readbody

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"webtest/colors"
)

type Eat struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
}

func ReadBody(r *http.Request, filePath string) error {
	filePath = filepath.Join(filePath, "data_files")

	err := deserialization(r, filePath)
	if err != nil {
		return err
	}

	return nil
}

func deserialization(r *http.Request, filePath string) error {
	var eat_data Eat

	//Чиатем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("Ошибка чтения тела запроса!")
	}
	colors.SetColor(colors.Text_Green)
	log.Println("Чтение тела запроса прошло успешно!")
	colors.ResetColor()

	//Процесс десериализации
	err = json.Unmarshal(body, &eat_data)
	if err != nil {
		return errors.New("Ошибка десериализации!")
	}
	colors.SetColor(colors.Text_Green)
	log.Println("Десериализация прошла успешно!")
	colors.ResetColor()

	colors.SetColor(colors.Text_Green)
	fmt.Println("\nResult: \n", eat_data)
	colors.ResetColor()

	filePath = filepath.Join(filePath, eat_data.Name+".json")
	err = serialization(&eat_data, &filePath)
	if err != nil {
		return err
	}

	return nil
}

func serialization(eat *Eat, filePath *string) error {
	//Процесс сериализации
	data, err := json.Marshal(eat)
	if err != nil {
		return errors.New("Ошибка сериализации!")
	}
	colors.SetColor(colors.Text_Green)
	log.Println("Сериализация прошла успешно!")
	colors.ResetColor()

	//Читаем каталог /data_files
	files, err := os.ReadDir(filepath.Dir(*filePath))
	if err != nil {
		return errors.New("Ошибка чтения каталога!")
	}

	for _, i := range files {
		if i.Name() == eat.Name+".json" {
			return errors.New("Блюдо уже существует!")
		}
	}

	//Записываем в новый файл
	err = os.WriteFile(*filePath, data, 0644)
	if err != nil {
		return errors.New("Ошибка записи в файл!")
	}
	colors.SetColor(colors.Text_Green)
	log.Println("Содержимое записано в файл!")
	colors.ResetColor()

	return nil
}
