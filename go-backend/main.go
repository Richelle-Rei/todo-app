package main

import (
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ToDoEntry struct{
	UserId int `json:"userId"`
	Id int `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

var data = []ToDoEntry{
    {1, 1, "delectus aut autem", false},
    {1, 2, "quis ut nam facilis et officia qui", false},
    {1, 3,  "fugiat veniam minus", false},
    {1, 4,  "et porro tempora", true},
	{1, 5,  "laboriosam mollitia et enim quasi adipisci quia provident illum", false},
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	
	e.GET("/todos", getToDoList)
	e.GET("/todos/:id", getToDo)
	e.POST("/todos", addToDo)
	e.PUT("/todos/:id", updateToDo)
	e.DELETE("/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(":3000"))

}

func getToDoList(c echo.Context) error{
	return c.JSON(http.StatusOK, data)
}

func getToDo(c echo.Context) error{
	id, err:= strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	for _, entry:= range data{
		if entry.Id == id{
			return c.JSON(http.StatusOK, entry)
		}
	}
	return c.JSON(http.StatusNotFound, "To Do not found")
}

func addToDo(c echo.Context) error{
	var newTodo ToDoEntry
	if err := c.Bind(&newTodo); err != nil{
		return c.JSON(http.StatusBadRequest, err)
	}
	id := data[len(data)-1].Id+1
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


func deleteTodo(c echo.Context) error{
	id, err:= strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID not found")
	}

	for index, entry:= range data{
		if entry.Id == id{
			data = append(data[:index], data[index+1:]...)
			
			return c.JSON(http.StatusOK, entry)
		}
	}


	return c.JSON(http.StatusNotFound, "To Do not found")
}


