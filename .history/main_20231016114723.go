package main

import (
	"errors"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func setupDatabase() (*sql.DB, error) {
    db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go_todo")
    if err != nil {
        return nil, err
    }
    return db, nil
}

func main() {
    db, err := setupDatabase()
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    // Rest of your code here...
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.POST("/todos", addTodo)
	router.Run("localhost:9000")
}


type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: false},
	{ID: "3", Item: "Record Video", Completed: false},
}

func getTodos(context *gin.Context) {
    rows, err := db.Query("SELECT id, item, completed FROM todos")
    if err != nil {
        context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve todos"})
        return
    }
    defer rows.Close()

    var todos []todo
    for rows.Next() {
        var t todo
        err := rows.Scan(&t.ID, &t.Item, &t.Completed)
        if err != nil {
            context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve todos"})
            return
        }
        todos = append(todos, t)
    }

    context.IndentedJSON(http.StatusOK, todos)
}

func addTodo(context *gin.Context) {
    var newTodo todo

    if err := context.BindJSON(&newTodo); err != nil {
        return
    }

    // Simpan data todo ke database
    _, err := db.Exec("INSERT INTO todos (item, completed) VALUES (?, ?)", newTodo.Item, newTodo.Completed)
    if err != nil {
        context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert todo"})
        return
    }

    context.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)

}

func getTodoById(id string) (*todo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("todo not found")
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	todo.Completed = !todo.Completed

	context.IndentedJSON(http.StatusOK, todo)
}

