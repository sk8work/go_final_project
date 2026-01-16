package db

import (
	"database/sql"
	"fmt"
)

// Migrate проверяет и обновляет схему БД если нужно
func Migrate(db *sql.DB) error {
	// Проверяем существование всех необходимых колонок
	columns := []string{"id", "date", "title", "comment", "repeat"}

	for _, column := range columns {
		if !columnExists(db, column) {
			return fmt.Errorf("column %s does not exist in table scheduler", column)
		}
	}

	return nil
}

// columnExists проверяет существование колонки в таблице
func columnExists(db *sql.DB, columnName string) bool {
	query := `
        SELECT COUNT(*) 
        FROM pragma_table_info('scheduler') 
        WHERE name = ?
    `

	var count int
	err := db.QueryRow(query, columnName).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}
