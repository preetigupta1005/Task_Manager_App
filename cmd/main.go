package main

import (
	"My-todo-app/database"
	"My-todo-app/server"
	"fmt"
	"log"
	"net/http"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	// db creds put in env and pass to the function

	err := database.ConnectDB()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.Close()
	r := server.SetUpRoutes()
	fmt.Println("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Server is not running", err)
	}
}
