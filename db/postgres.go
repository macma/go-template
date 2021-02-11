package db

import (
	"context"
	"database/sql"
	"eass/bpe-scheduler/config"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var conn *sql.DB

// InitDB initializes database connection with provided connection
func InitDBWithConn(newConn *sql.DB) {
	conn = newConn
}

// InitDB initializes database connection
func InitDB() (db *sql.DB) {
	connStr := `postgres://` + config.Config.PostgresUsername +
		`:` + config.Config.PostgresPassword +
		`@` + config.Config.PostgresAddress +
		`/` + config.Config.PostgresDatabase +
		`?sslmode=` + config.Config.PostgresSslMode

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		// logger.Panicw(ctx, fmt.Sprintf("there's error during make connection to db. %v", err))
	}

	if err = db.Ping(); err != nil {
		// logger.Panicw(ctx, fmt.Sprintf("there's error during ping connection to db. %v", err))
	}

	conn = db

	fmt.Println("Initialized Database successfully")
	// logger.Infow(ctx, fmt.Sprint("Initialized Database successfully"))
	MigrateDB()
	return db
}

func MigrateDB() {
	fmt.Println(fmt.Sprint("Start with migrate DB"))

	driver, err := postgres.WithInstance(conn, &postgres.Config{MigrationsTable: "schema_migrations"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Get migrate db driver information error: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance("file://./db/migrations", "postgres", driver)
	if err == nil {
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			fmt.Println(fmt.Sprintf("Migrate db UP error: %v", err))
		}
		fmt.Println(fmt.Sprintf("Migrate DB successful!, the migration path is: %s", "file://./db/migrations"))
	} else {
		fmt.Println(fmt.Sprintf("migrate.NewWithDatabaseInstance fail! %s", err))
	}
}

func getSliceFromRow(ctx context.Context, rows *sql.Rows) (cols []string, vals [][]string) {

	cols, err := rows.Columns()
	if err != nil {
		// logger.Errorw(ctx, fmt.Sprintf("Failed to get columns, %v", err))
		return
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	var resultArray [][]string

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make([]string, len(cols))
		err = rows.Scan(dest...)
		if err != nil {
			// logger.Errorw(ctx, fmt.Sprintf("Failed to scan row, %v", err))
			return
		}
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
		resultArray = append(resultArray, result)
	}

	return cols, resultArray
}
