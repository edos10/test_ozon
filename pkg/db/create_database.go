package databases

import "database/sql"

func checkTableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`

	var exists bool
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func createTableForString(db *sql.DB) error {
	query := `
		CREATE TABLE genstring (
			currentstring TEXT
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func createTableForUrls(db *sql.DB) error {
	query := `
		CREATE TABLE urls (
			original_url TEXT,
			short_url TEXT
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
