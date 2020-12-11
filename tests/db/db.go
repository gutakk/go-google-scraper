package tests

import "fmt"

func ConstructTestDsn() string {
	host := "localhost"
	port := "5433"
	dbName := "go_google_scraper_test"
	username := "postgres"

	return fmt.Sprintf("sslmode=disable host=%s port=%s dbname=%s user=%s",
		host,
		port,
		dbName,
		username,
	)
}
