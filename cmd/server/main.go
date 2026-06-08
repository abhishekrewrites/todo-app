package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"todo-app/internals/database"
	"todo-app/internals/handlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system env vars")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database connected successfully")

	// -----****Auth****-----

	authHandler := handlers.NewAuthHandler(db)

	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	// -----****Todos****-----

	todoHandler := handlers.NewTodoHandler(db)

	http.HandleFunc("/todos", todoHandler.GetTodos)
	http.HandleFunc("/todos/create", todoHandler.CreateTodo)
	http.HandleFunc("/todos/update", todoHandler.UpdateTodo)
	http.HandleFunc("/todos/delete", todoHandler.DeleteTodo)

	log.Println("Server started on :8080")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server started on :" + port)

	http.ListenAndServe(":"+port, nil)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
