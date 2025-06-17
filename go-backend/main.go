package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type ToDoEntry struct {
	UserId    int    `json:"userId" db:"user_id"`
	Id        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Completed bool   `json:"completed" db:"completed"`
}

var data = []ToDoEntry{
	{1, 1, "delectus aut autem", false},
	{1, 2, "quis ut nam facilis et officia qui", false},
	{1, 3, "fugiat veniam minus", false},
	{1, 4, "et porro tempora", true},
	{1, 5, "laboriosam mollitia et enim quasi adipisci quia provident illum", false},
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())


	connStr := "user=richellereivanka dbname=tododb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	e.GET("/todos", func(c echo.Context) error { return getToDoSQL(c, db) })
	e.POST("/todos", func(c echo.Context) error { return addToDoSQL(c, db) })
	e.PUT("/todos/:id", func(c echo.Context) error { return UpdateToDoSQL(c, db) })
	e.DELETE("/todos/:id", func(c echo.Context) error { return deleteToDoSQL(c, db) })

	e.GET("/todos/:id", getToDo)
	// e.GET("/todos", getToDoList)
	// e.POST("/todos", addToDo)
	// e.PUT("/todos/:id", updateToDo)
	// e.DELETE("/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(":3000"))

}

func getToDoSQL(c echo.Context, db *sql.DB) error {
	var todoList []ToDoEntry
	rows, err := db.Query("SELECT user_id, id, title, completed FROM todo ORDER BY id")

	if err != nil {
		log.Println("Database Query Error", err)
		fmt.Println("Database Query Error")
		return c.JSON(http.StatusBadRequest, "Database Query Error")
	}

	defer rows.Close()

	for rows.Next() {
		var each ToDoEntry
		err := rows.Scan(&each.UserId, &each.Id, &each.Title, &each.Completed)
		if err != nil {
			log.Println("Database Row Scan Error", err)
			fmt.Println("Database Row Scan Error")
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

	_, err := db.Exec("INSERT INTO todo (user_id, title, completed) VALUES ($1, $2, $3);", newTodo.UserId, newTodo.Title, newTodo.Completed)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, "Insert to Database Error")
	}
	return c.JSON(http.StatusCreated, newTodo)

}

func UpdateToDoSQL(c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	var updatedData struct {
		Completed bool `json:"completed" db:"completed`
	}

	if err := c.Bind(&updatedData); err != nil {
		return c.JSON(http.StatusBadRequest, "Completed not found")
	}
	_, err = db.Exec("UPDATE todo SET completed = $1 WHERE id = $2", updatedData.Completed, id)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, "Delete from Database Error")
	}

	return getToDoSQL(c, db)

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

	return getToDoSQL(c, db)

}

func getToDoList(c echo.Context) error {
	return c.JSON(http.StatusOK, data)
}

func getToDo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	for _, entry := range data {
		if entry.Id == id {
			return c.JSON(http.StatusOK, entry)
		}
	}
	return c.JSON(http.StatusNotFound, "To Do not found")
}

func addToDo(c echo.Context) error {
	var newTodo ToDoEntry
	if err := c.Bind(&newTodo); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	id := data[len(data)-1].Id + 1
	newTodo.Id = id
	data = append(data, newTodo)

	return c.JSON(http.StatusCreated, newTodo)
}

func updateToDo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	for index, entry := range data {
		if entry.Id == id {
			data[index].Completed = !data[index].Completed
			return c.JSON(http.StatusOK, data[index])
		}
	}
	return c.JSON(http.StatusNotFound, "To Do not found")
}

func deleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	for index, entry := range data {
		if entry.Id == id {
			data = append(data[:index], data[index+1:]...)

			return c.JSON(http.StatusOK, entry)
		}
	}

	return c.JSON(http.StatusNotFound, "To Do not found")
}
