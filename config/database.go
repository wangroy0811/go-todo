package config

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/thoth"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

func Database() *sql.DB {
	logger, _ := thoth.Init("log")

	user, exist := os.LookupEnv("MYSQL_USER")

	if !exist {
		logger.Log(errors.New("MYSQL_USER not set in .env"))
		log.Fatal("MYSQL_USER not set in .env")
	}

	pass, exist := os.LookupEnv("MYSQL_PASSWORD")

	if !exist {
		logger.Log(errors.New("MYSQL_PASSWORD not set in .env"))
		log.Fatal("MYSQL_PASSWORD not set in .env")
	}

	host, exist := os.LookupEnv("MYSQL_HOST")

	if !exist {
		logger.Log(errors.New("MYSQL_HOST not set in .env"))
		log.Fatal("MYSQL_HOST not set in .env")
	}

	port, exist := os.LookupEnv("MYSQL_PORT")

	if !exist {
		logger.Log(errors.New("MYSQL_PORT not set in .env"))
		log.Fatal("MYSQL_PORT not set in .env")
	}

	credentials := fmt.Sprintf("%s:%s@(%s:%s)/?charset=utf8&parseTime=True", user, pass, host, port)

	database, err := sql.Open("mysql", credentials)

	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	} else {
		fmt.Println("Database Connection Successful")
	}

	_, err = database.Exec(`CREATE DATABASE gotodo`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Exec(`USE gotodo`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Exec(`
		CREATE TABLE todos (
		    id INT AUTO_INCREMENT,
		    item TEXT NOT NULL,
		    completed BOOLEAN DEFAULT FALSE,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    PRIMARY KEY (id)
		);
	`)

	if err != nil {
		fmt.Println(err)
	}

	// Check if created_at column exists
	var count int
	err = database.QueryRow(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = 'gotodo' 
		AND table_name = 'todos' 
		AND column_name = 'created_at'
	`).Scan(&count)

	if err != nil {
		fmt.Println("Error checking if created_at column exists:", err)
	}

	// Add created_at column if it doesn't exist
	if count == 0 {
		_, err = database.Exec(`ALTER TABLE todos ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)
		if err != nil {
			fmt.Println("Error adding created_at column:", err)
		}
	}

	return database
}