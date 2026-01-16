package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// SQL схема для создания таблицы
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '' CHECK(length(date) = 8 OR date = ''),
    title VARCHAR(255) NOT NULL DEFAULT '' CHECK(length(title) <= 255),
    comment TEXT CHECK(length(comment) <= 1000),
    repeat VARCHAR(128) DEFAULT '' CHECK(length(repeat) <= 128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date 
ON scheduler(date);
`

// Init инициализирует базу данных
func Init(dbFile string) error {
	// Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем соединение с БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Если файл не существовал, создаём таблицу
	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return fmt.Errorf("failed to create schema: %w", err)
		}
		fmt.Printf("Database created at: %s\n", dbFile)
	} else {
		// Проверяем, существует ли таблица scheduler
		var tableExists int
		checkTableSQL := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='scheduler'`
		if err := db.QueryRow(checkTableSQL).Scan(&tableExists); err != nil {
			db.Close()
			return fmt.Errorf("failed to check table existence: %w", err)
		}

		if tableExists == 0 {
			// Таблицы нет, создаём её
			if _, err := db.Exec(schema); err != nil {
				db.Close()
				return fmt.Errorf("failed to create schema: %w", err)
			}
			fmt.Printf("Table created in existing database: %s\n", dbFile)
		}
	}

	DB = db
	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB возвращает экземпляр базы данных
func GetDB() *sql.DB {
	return DB
}

// IsConnected проверяет, подключена ли база данных
func IsConnected() bool {
	if DB == nil {
		return false
	}

	// Проверяем пинг
	if err := DB.Ping(); err != nil {
		return false
	}

	return true
}
