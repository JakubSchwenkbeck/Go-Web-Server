package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver package for database interaction.
)

var db *sql.DB // Global variable to hold the database connection pool.

/*
InitDB initializes the database connection using the provided data source name.

This function sets up a connection pool to the MySQL database. The dataSourceName
should include the necessary connection information, such as the username, password,
host, and database name.

Example:

	dataSourceName := "username:password@tcp(127.0.0.1:3306)/dbname"

Parameters:
  - dataSourceName: A string containing the connection details to the MySQL database.

Behavior:
  - Initializes the global `db` variable with a connection pool.
  - Logs and exits the program if the connection cannot be established.
  - Verifies the connection by pinging the database.
*/
func InitDB(dataSourceName string) {
	var err error

	// Open a connection to the database using the provided data source name.
	// sql.Open does not actually establish any connections but prepares the database connection for future use.
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		// Log the error and exit the program if the connection cannot be established.
		log.Fatal(err)
	}

	// Ping the database to verify the connection is valid.
	// Ping establishes a connection to the database and verifies that it can be communicated with.
	if err = db.Ping(); err != nil {
		// Log the error and exit the program if the connection cannot be verified.
		log.Fatal(err)
	}

	// Print a success message to indicate that the database connection is established.
	fmt.Println("Connected to the database!")
}

/*
connectDB demonstrates how to establish a database connection locally within a function.

This function uses a connection string to connect to the MySQL database. The connection
string includes the username, password, host, and database name.

Behavior:
  - Opens a new database connection using a locally defined connection string.
  - Defers the closing of the database connection to ensure it is closed when the function exits.
  - Verifies the connection by pinging the database.
  - Prints a success message if the connection is successful.
  - Panics if the connection cannot be established or verified.
*/
func connectDB() {

	// Define the connection string with the necessary details to connect to the database.
	// Example format: "username:password@tcp(127.0.0.1:3306)/dbname"
	db, err := sql.Open("mysql", getDBString())
	if err != nil {
		// Panic is used to stop the program execution and print the error message if the connection fails.
		panic(err)
	}
	// Ensure the database connection is closed when the function exits to avoid connection leaks.
	defer db.Close()

	// Ping the database to verify that the connection is successful.
	err = db.Ping()
	if err != nil {
		// Panic is used to stop the program execution and print the error message if the connection verification fails.
		panic(err)
	}

	// Print a success message to indicate that the database connection is established.
	fmt.Println("Successfully connected to MySQL database!")
}
