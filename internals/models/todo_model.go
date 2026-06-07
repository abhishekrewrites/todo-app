package models

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}


type CreateTodoRequest struct {
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}


type UpdateTodoRequest struct {
	ID        int  `json:"id"`
	Completed bool `json:"completed"`

}


type DeleteTodoRequest struct {
	ID        int  `json:"id"`
}
