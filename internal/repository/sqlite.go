package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/paudarco/todo/internal/config"
	_ "modernc.org/sqlite"
)

// NewSQLiteDB запускает в работу базу данных.
func NewSQLiteDB(cfg config.DB) (*sql.DB, error) {
	// Получаем файл базы данных, если такой существует
	// или есть возможность его создать.
	dbFile, err := CreateDBFile(cfg)
	if err != nil {
		return nil, err
	}

	// Открываем подключение к базе данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// Запускаем миграции.
	err = InitTable(db)
	if err != nil {
		return nil, err
	}

	// err = TestDB(db)
	// if err != nil {
	// 	return nil, err
	// }

	return db, nil
}

// CreateDBFile создает или возвращает существующий файл базы данных.
func CreateDBFile(cfg config.DB) (string, error) {
	dbFile := cfg.Path

	// Проверяем файл бд на существование.
	// Если его не существует - создаем.
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			return "", err
		}
		err = file.Close()
		if err != nil {
			return "", err
		}
	}
	return dbFile, nil
}

// InitTable создает таблицу с задачами и индекс.
func InitTable(db *sql.DB) error {
	// Открываем транзкцию для создания таблицы и индекса.
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date TEXT NOT NULL,
				title TEXT NOT NULL,
				comment TEXT,
				repeat TEXT
			);`

	_, err = tx.Exec(query)
	if err != nil {
		return tx.Rollback()
	}

	query = "CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);"
	_, err = tx.Exec(query)
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

// TestDB тестирует работу базы данных.
func TestDB(db *sql.DB) error {
	var id int
	row := db.QueryRow(`INSERT INTO scheduler (date, title, comment, repeat) VALUES ("2022-5-7", $1, $2, $3) returning id`, "dada", "bebe", "dadada")
	if err := row.Scan(&id); err != nil {
		return err
	}
	fmt.Println("id: ", id)
	return nil
}
