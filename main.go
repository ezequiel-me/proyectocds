package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Tasks struct {
	ID          int       `json: "id"`
	Title       string    `json: "title"`
	Description string    `json: "description"`
	Status      string    `json: "status"`
	CreatedAt   time.Time `json: "createdAt"`
	CompletedAt time.Time `json: "completedAt"`
}

var TasksDB = []Tasks{}

func main() {
	//inicializamos las rutas
	ru := gin.Default()

	//Si se recibe peticion por metodo get por esa ruta, se envia el handler de getTasks
	ru.GET("/api/tasks", getTasks)
	ru.POST("/api/tasks", postTasks)
	ru.GET("/api/tasks/:id", getTaskById)
	ru.PUT("/api/tasks/:id", updateTaskById)
	ru.DELETE("/api/tasks/:id", deleteTaskById)
	ru.GET("/api/tasks/search", getTaskByName)

	//se ejecuta el servidor
	ru.Run()
}
