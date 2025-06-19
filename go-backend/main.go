package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type ToDoEntry struct {
	UserId       int    `json:"userId" db:"user_id"`
	Id           int    `json:"id" db:"id"`
	Title        string `json:"title" db:"title"`
	Completed    bool   `json:"completed" db:"completed"`
	DisplayOrder int    `json:"displayOrder" db:"display_order"`
}

// var data = []ToDoEntry{
// 	{1, 1, "delectus aut autem", false},
// 	{1, 2, "quis ut nam facilis et officia qui", false},
// 	{1, 3, "fugiat veniam minus", false},
// 	{1, 4, "et porro tempora", true},
// 	{1, 5, "laboriosam mollitia et enim quasi adipisci quia provident illum", false},
// }

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	godotenv.Load(".env")

	m, err := migrate.New(
		"file://db/migrations",
		os.Getenv("SQL_URL"),
	)
	if err != nil {
		log.Fatal("Error in setting up migration", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error in running migration up", err)
	}
	defer m.Close()

	connStr := os.Getenv("CONNECT_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e.GET("/todos", func(c echo.Context) error { return getToDoListSQL(c, db) })
	e.POST("/todos", func(c echo.Context) error { return addToDoSQL(c, db) })
	e.PUT("/todos/:id", func(c echo.Context) error { return updateToDoSQL(c, db) })
	e.PUT("/todos/:id/edit", func(c echo.Context) error { return editToDoSQL(c, db) })
	e.PUT("/todos/order", func(c echo.Context) error { return updateDisplaySQL(c, db) })
	e.DELETE("/todos/:id", func(c echo.Context) error { return deleteToDoSQL(c, db) })

	// e.GET("/todos/:id", getToDo)
	// e.GET("/todos", getToDoList)
	// e.POST("/todos", addToDo)
	// e.PUT("/todos/:id", updateToDo)
	// e.DELETE("/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(":3000"))

}

func getToDoListSQL(c echo.Context, db *sql.DB) error {
	var todoList []ToDoEntry
	rows, err := db.Query("SELECT user_id, id, title, completed, display_order FROM todo ORDER BY display_order")

	if err != nil {
		log.Println("Database Query Error", err)
		return c.JSON(http.StatusBadRequest, "Database Query Error")
	}

	defer rows.Close()

	for rows.Next() {
		var each ToDoEntry
		err := rows.Scan(&each.UserId, &each.Id, &each.Title, &each.Completed, &each.DisplayOrder)
		if err != nil {
			log.Println("Database Row Scan Error", err)
			return c.JSON(http.StatusBadRequest, "Database Row Scan Error")
		}

		todoList = append(todoList, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	for _, each := range todoList {
		fmt.Println(each.Id, each.Title)
	}
	return c.JSON(http.StatusOK, todoList)

}

func addToDoSQL(c echo.Context, db *sql.DB) error {
	var newTodo ToDoEntry
	if err := c.Bind(&newTodo); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err := db.Exec("INSERT INTO todo (user_id, title, completed, display_order) VALUES ($1, $2, $3, $4);", newTodo.UserId, newTodo.Title, newTodo.Completed, newTodo.DisplayOrder)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, "Insert to Database Error")
	}
	return c.JSON(http.StatusCreated, newTodo)

}

func updateToDoSQL(c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	var updatedData struct {
		Completed bool `json:"completed" db:"completed"`
	}

	if err := c.Bind(&updatedData); err != nil {
		return c.JSON(http.StatusBadRequest, "Completed not found")
	}
	_, err = db.Exec("UPDATE todo SET completed = $1 WHERE id = $2", updatedData.Completed, id)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, "Update Completed Error")
	}

	return getToDoListSQL(c, db)

}

func editToDoSQL(c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	var updatedData struct {
		Title string `json:"title" db:"title"`
	}

	if err := c.Bind(&updatedData); err != nil {
		return c.JSON(http.StatusBadRequest, "Edited title not found")
	}

	if updatedData.Title == "" {
		return c.JSON(http.StatusBadRequest, "Title cannot be empty")
	}

	_, err = db.Exec("UPDATE todo SET title = $1 WHERE id = $2", updatedData.Title, id)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, "Update Completed Error")
	}

	return getToDoListSQL(c, db)

}

func deleteToDoSQL(c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}
	_, err = db.Exec("DELETE FROM todo WHERE id = $1", id)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, "Delete from Database Error")
	}

	return getToDoListSQL(c, db)

}

func updateDisplaySQL(c echo.Context, db *sql.DB) error {
	var updatedData struct {
		NewList []ToDoEntry `json:"newList"`
	}
	if err := c.Bind(&updatedData); err != nil {
		fmt.Print("ToDo Data not found")
		return c.JSON(http.StatusBadRequest, "ToDo Data not found")
	}

	_, err := db.Exec("UPDATE todo SET display_order = -id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Reset Order Failed")
	}

	for _, each := range updatedData.NewList {
		_, err := db.Exec("UPDATE todo SET display_order = $1 WHERE id = $2", each.DisplayOrder, each.Id)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusBadRequest, "Updata Display Order Error")
		}
	}
	return getToDoListSQL(c, db)
}

// func getToDoList(c echo.Context) error {
// 	return c.JSON(http.StatusOK, data)
// }

// func getToDo(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))

// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "ID not found")
// 	}

// 	for _, entry := range data {
// 		if entry.Id == id {
// 			return c.JSON(http.StatusOK, entry)
// 		}
// 	}
// 	return c.JSON(http.StatusNotFound, "To Do not found")
// }

// func addToDo(c echo.Context) error {
// 	var newTodo ToDoEntry
// 	if err := c.Bind(&newTodo); err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	id := data[len(data)-1].Id + 1
// 	newTodo.Id = id
// 	data = append(data, newTodo)

// 	return c.JSON(http.StatusCreated, newTodo)
// }

// func updateToDo(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "ID not found")
// 	}

// 	for index, entry := range data {
// 		if entry.Id == id {
// 			data[index].Completed = !data[index].Completed
// 			return c.JSON(http.StatusOK, data[index])
// 		}
// 	}
// 	return c.JSON(http.StatusNotFound, "To Do not found")
// }

// func deleteTodo(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "ID not found")
// 	}

// 	for index, entry := range data {
// 		if entry.Id == id {
// 			data = append(data[:index], data[index+1:]...)

// 			return c.JSON(http.StatusOK, entry)
// 		}
// 	}

// 	return c.JSON(http.StatusNotFound, "To Do not found")
// }
