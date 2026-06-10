package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"webtest/internal/logger"
	"webtest/internal/models"
	"webtest/pkg/colors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	dataBase   *sql.DB
	driverName string
)

// Errors
var (
	ConnectionErr   = errors.New("Ошибка подключения к БД!")
	GetVersionErr   = errors.New("Не удалось получить версию БД!")
	CreateTableErr  = errors.New("Ошибка создания таблицы!")
	RecipeExistsErr = errors.New("Рецепт с таким названием уже существует!")
	GetDataErr      = errors.New("Не удалось получить результат!")
	ReadDataErr     = errors.New("Ошибка чтения результата!")
	DeleteDataErr   = errors.New("Не удалось выполнить удаление!")
	UpdateDataErr   = errors.New("Не удалось выполнить обновление!")
	ProcessingErr   = errors.New("Ошибка в обработке результатов!")
	ctxErr          = errors.New("Пользователь отключился! [DATABASE]")
)

// Несколько попыток подключения к БД
func ConnectionAttempt(dsn string) error {
	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		log.Printf("Попытка подключения к PostgreSQL №%d...\t", i)

		db, err = sql.Open(driverName, dsn)

		if err == nil {
			break
		}

		log.Printf("Ошибка подключения: \"%v\"\n", err.Error())
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return ConnectionErr
	}

	dataBase = db
	logger.LogPrint("Подключение к БД выполнено успешно!", models.Info)
	return nil
}

// Подключение к БД
func ConnectDB() (*sql.DB, string, error) {
	//Получение значений из переменных окружения ОС
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	driverName = "postgres"

	//Подключение к БД
	err := ConnectionAttempt(dsn)
	if err != nil {
		return nil, "", err
	}

	//Получение версии БД
	rows, err := dataBase.Query(`
		select version()
		`)
	if err != nil {
		return nil, "", GetVersionErr
	}

	var version string
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			colors.SetColor(colors.Text_Red)
			logger.LogPrint(
				fmt.Sprintf("Не удалось получить результат запроса: %v", err.Error()),
				models.Error,
			)
			colors.ResetColor()
		}
	}
	if err = rows.Err(); err != nil {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(
			fmt.Sprintf("Ошибка при обработке результатов: %v", err.Error()),
			models.Error,
		)
		colors.ResetColor()
	}

	// Создаем таблицу recipes
	err = createTables()
	if err != nil {
		return nil, "", CreateTableErr
	}

	return dataBase, version, nil
}

// Создание таблицы
func createTables() error {
	//Таблица recipes
	_, err := dataBase.Query(`
	create table if not exists recipes (
	id serial primary key,
	name text,
	category text,
	ingredients text[],
	instructions text
	)
	`)

	return err
}

// Запись данных
func WriteDB(eat_data *models.Eat) error {
	//Проверка на наличие повторных данных
	var exists bool = false
	err := dataBase.QueryRow(`
	select exists (
	select 1 from recipes where name ILIKE $1
	)
	`, eat_data.Name).Scan(&exists)

	if err != nil {
		return GetDataErr
	}

	if exists {
		return RecipeExistsErr
	}

	//Запись данных в БД
	dataBase.QueryRow(`
	insert into recipes (
	name,
	category,
	ingredients,
	instructions
	)
	values (
	$1,
	$2,
	$3,
	$4)
	`, eat_data.Name, eat_data.Category, pq.Array(eat_data.Ingredients), eat_data.Instructions)

	return nil
}

// Получение всех рецептов
func GetAllRecipes(ctx context.Context) ([]models.RecipeShort, error) {
	var data []models.RecipeShort

	rows, err := dataBase.QueryContext(ctx, `
	select id, name from recipes
	`)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ctxErr
		}

		return nil, GetDataErr
	}

	for rows.Next() {
		var recipe models.RecipeShort
		err = rows.Scan(
			&recipe.Id,
			&recipe.Name,
		)
		if err != nil {
			return nil, ReadDataErr
		}
		data = append(data, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, ProcessingErr
	}

	return data, nil
}

// Получение одного рецепта
func GetOneRecipe(id int) (models.Eat, error) {
	var data models.Eat
	err := dataBase.QueryRow(`
	select name, category, ingredients, instructions
	from recipes
	where id = $1
	`, id).Scan(
		&data.Name,
		&data.Category,
		pq.Array(&data.Ingredients),
		&data.Instructions,
	)
	if err != nil {
		return data, GetDataErr
	}

	return data, nil
}

// Удаление рецепта
func DeleteRecipe(id int) error {
	_, err := dataBase.Query(`
	delete from recipes
	where id = $1
	`, id)
	if err != nil {
		return DeleteDataErr
	}

	return nil
}

// Обновление рецепта
func UpdateRecipe(eat_data *models.Eat) error {
	_, err := dataBase.Query(`
	update recipes
	set
	name = $1,
	category = $2,
	ingredients = $3,
	instructions = $4
	where id = $5
	`,
		&eat_data.Name,
		&eat_data.Category,
		pq.Array(&eat_data.Ingredients),
		&eat_data.Instructions,
		&eat_data.ID,
	)
	if err != nil {
		return UpdateDataErr
	}

	return nil
}
