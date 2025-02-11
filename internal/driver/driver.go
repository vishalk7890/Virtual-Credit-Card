package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// func OpenDB(dsn string) (*sql.DB, error) {
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

//		return db, nil
//	}
func OpenDB(dsn string) (*sql.DB, error) {
	fmt.Println("Connecting to database with DSN:", dsn) // Debugging output

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return nil, err
	}

	fmt.Println("Database connection successful")
	return db, nil
}
