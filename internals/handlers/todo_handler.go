package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"todo-app/internals/models"
)

type TodoHandler struct {
	DB *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{
		DB: db,
	}
}

func (h *TodoHandler) GetTodos(
	w http.ResponseWriter,
	r *http.Request,
) {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset := (page - 1) * limit

	rows, err := h.DB.Query(
		`
		SELECT id, title, completed, created_at
		FROM todos
		ORDER BY id
		LIMIT $1 OFFSET $2
		`,
		limit,
		offset,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	var todos []models.Todo

	for rows.Next() {

		var todo models.Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	var total int

	err = h.DB.QueryRow(
		`
		SELECT COUNT(*)
		FROM todos
		`,
	).Scan(&total)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
			"data":  todos,
		},
	)
}

func (h *TodoHandler) CreateTodo(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req models.CreateTodoRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	var id int

	err = h.DB.QueryRow(
		`
		INSERT INTO todos(title)
		VALUES($1)
		RETURNING id
		`,
		req.Title,
	).Scan(&id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"id":      id,
			"message": "Todo created successfully",
		},
	)
}



func (h *TodoHandler) UpdateTodo(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req models.UpdateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	result, err := h.DB.Exec(
		`
		UPDATE todos
		SET completed = $1
		WHERE id = $2
		`,
		req.Completed,
		req.ID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"rowsAffected": rowsAffected,
			"message":      "Todo updated",
		},
	)
}


func (h *TodoHandler) DeleteTodo(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req models.DeleteTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	result, err := h.DB.Exec(
		`
		DELETE FROM todos
		WHERE id = $1
		`,
		req.ID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"rowsAffected": rowsAffected,
			"message":      "Todo deleted successfully",
		},
	)
}


