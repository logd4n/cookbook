package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"webtest/internal/models"
	"webtest/pkg/colors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	dataBase   *sql.DB
	driverName string
)

func ConnectDB() (*sql.DB, string) {
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
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		colors.SetColor(colors.Text_Red)
		log.Fatal("Ошибка подключения к БД: ", err)
		colors.ResetColor()
	}
	dataBase = db

	//Получение версии БД
	rows, err := db.Query(`
		select version()
		`)
	if err != nil {
		colors.SetColor(colors.Text_Red)
		log.Fatal("Не удалось получить версию БД!", err)
		colors.ResetColor()
	}

	var version string
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			colors.SetColor(colors.Text_Red)
			log.Fatal("Ошибка чтения результата! (DB version)", err)
			colors.ResetColor()
		}
	}
	if err = rows.Err(); err != nil {
		colors.SetColor(colors.Text_Red)
		log.Fatal("Ошибка при обработке результатов! (DB version)", err)
		colors.ResetColor()
	}

	// Создаем таблицу recipes
	err = createTables()
	if err != nil {
		colors.SetColor(colors.Text_Red)
		log.Fatal(err.Error())
	}

	colors.ResetColor()

	return db, version
}

func createTables() error {
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

func WriteDB(eat_data *models.Eat) error {
	//Проверка на наличие повторных данных
	var exists bool = false
	err := dataBase.QueryRow(`
	select exists (
	select 1 from recipes where name ILIKE $1
	)
	`, eat_data.Name).Scan(&exists)

	if err != nil {
		colors.SetColor(colors.Text_Red)
		fmt.Println("Не удалось проверить наличие повторных данных")
		colors.ResetColor()
	}

	if exists {
		return errors.New("Рецепт с таким названием уже существует!")
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

	/*
		rows, err := dataBase.Query(`
		select *
		from recipes
		where name = $1
		`)
		if err != nil {
			colors.SetColor(colors.Text_Red)
			log.Fatal("Не удалось выполнить запрос!", err)
			colors.ResetColor()
		}

		var scanName string
		for rows.Next() {
			err = rows.Scan(&scanName)
			if err != nil {
				colors.SetColor(colors.Text_Red)
				log.Fatal("Ошибка чтения результата!", err)
				colors.ResetColor()
			}
		}*/

	return nil
}

func GetAllRecipes() ([]models.RecipeShort, error) {
	var data []models.RecipeShort

	rows, err := dataBase.Query(`
	select id, name from recipes
	`)
	if err != nil {
		return nil, errors.New("Не удалось получить данные!")
	}

	for rows.Next() {
		var recipe models.RecipeShort
		err = rows.Scan(
			&recipe.Id,
			&recipe.Name,
		)
		if err != nil {
			return nil, errors.New("Ошибка чтения результата!")
		}
		data = append(data, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("Ошибка в обработке результатов!")
	}

	return data, nil
}

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
		return data, errors.New("Не удалось получить данные!")
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
		return errors.New("Не удалось выполнить удаление!")
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
		return errors.New("Не удалось выполнить обновление!")
	}

	return nil
}
