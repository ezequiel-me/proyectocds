package main

//HACER MUCHAS PRUEBAS PARA QUE NO SE PRODUZCAN ERRORES
/*
ERRORES DETECTADOS:
	LOS ID SE REPITEN, AL ELIMINAR UNA TAREA Y LUEGO AL COMENZAR A GENERAR MAS TAREAS.
*/
import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// handler para recibir la peticion del cliente
func getTasks(c *gin.Context) {
	//devolvemos el slice en formato json
	stat := c.Query("status")
	if strings.Compare(stat, "") != 0 {
		var taskStatus = []Tasks{}
		for _, sta := range TasksDB {
			if strings.Compare(stat, sta.Status) == 0 {
				taskStatus = append(taskStatus, sta)
			}
		}
		if len(taskStatus) <= 0 {
			c.IndentedJSON(http.StatusOK, 200)
			return
		}
		c.IndentedJSON(http.StatusOK, taskStatus)
		return
	} else {
		if len(TasksDB) <= 0 {
			c.IndentedJSON(http.StatusAccepted, gin.H{
				"Error message": "There isnt nothing in the slice",
				"QuantyTask":    len(TasksDB),
			},
			)
			return
		}
		c.IndentedJSON(http.StatusAccepted, TasksDB)
	}

}
func postTasks(c *gin.Context) {
	//variable de tipo Task para almecenar la data del peticion del cliente
	var task Tasks
	var newTask = []Tasks{}
	createdAt := time.Now()

	//Se envia la refencia de la variable "task" para poder almacenar la data de la peticion (el body) del cliente
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(http.StatusBadRequest, 400)
	}

	//obtenemos el ultimo elemento de la lista para generar el ID (ver si lo cambio por uno auto o lo hago en otra funcion asi solo es llamarla)
	numId := len(TasksDB)
	if numId <= 0 {
		numId++
	} else {
		numId = len(TasksDB) + 1
	}

	statusLower := strings.ToLower(task.Status)
	if statusLower != "new" {
		c.IndentedJSON(http.StatusBadRequest, &statusLower)
		return
	}
	//ver si es posible hacer el append de otra forma mas sencilla
	newTask = append(newTask, Tasks{ID: numId, Title: task.Title, Description: task.Description, Status: statusLower, CreatedAt: createdAt, CompletedAt: task.CompletedAt})

	TasksDB = append(TasksDB, newTask...)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	for _, t := range TasksDB {
		if t.ID == aid {
			c.IndentedJSON(http.StatusOK, t)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "element not found"})
}

func updateTaskById(c *gin.Context) {
	/*VER SI PONGO QUE SOLO SE PUEDAN ACTUALIZAR UNOS CAMPOS O TODOS!!!
	porque pide lo siguiente:
	Devuelve Status 400 Bad Request en el caso de que alguno de los campos necesarios para crear la task no vengan o el estado no sea uno de los permitidos.
	*/
	var upTask = []Tasks{}
	var task Tasks
	id, err := strconv.Atoi(c.Param("id"))
	//incrementar para cuando solo haya una tarea
	if err != nil {
		return
	}
	if id > len(TasksDB) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error with id"})
		return
	}

	if err := c.BindJSON(&task); err != nil {
		return
	}

	statusLower := strings.ToLower(task.Status)
	var completedAt time.Time
	if strings.Compare(statusLower, "completed") == 0 {
		completedAt = time.Now()
	} else if statusLower != "ongoing" && statusLower != "new" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "status not valid"})
		return
	}

	//INTENTAR HACERLO CON UN PUNTERO A MAYORES, ASI NOS AHORRAMOS EL TENER QUE CREAR OTRA VARIABLE

	var taskNew = TasksDB[id-1]
	upTask = append(upTask, Tasks{ID: id, Title: task.Title, Description: task.Description, Status: statusLower, CreatedAt: taskNew.CreatedAt, CompletedAt: completedAt})

	//actualizamos la posicion correcta del slice, por el contenido del slice Uptask en la misma posicion
	TasksDB[id-1] = upTask[0]
	c.IndentedJSON(http.StatusCreated, upTask)
}

func deleteTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id > len(TasksDB) {
		c.IndentedJSON(http.StatusOK, 200)
		return
	}
	var newTak []Tasks
	for _, t := range TasksDB {
		if t.ID != id {
			newTak = append(newTak, t)
		}
	}
	TasksDB = newTak
	c.IndentedJSON(http.StatusAccepted, gin.H{
		"message": "task removed succesfully",
		"id":      id,
	})
}

// INTENTAR HACERLO CON UN ARRAY
func getTaskByName(c *gin.Context) {
	title := c.Query("title")
	var tasks = []Tasks{}
	for _, tit := range TasksDB {
		if tit.Title == title {
			tasks = append(tasks, tit)
		}
	}
	if len(tasks) <= 0 {
		c.IndentedJSON(http.StatusOK, 200)
		return
	}
	c.IndentedJSON(http.StatusAccepted, tasks)
}
