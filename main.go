package main

import (
	elogrus "github.com/dictor/echologrus"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	Tasks map[int]*Task = map[int]*Task{}
	TaskId int = 0
)

type CustomValidator struct {
		validator *validator.Validate
	}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	elogrus.Attach(e)
	e.Validator = &CustomValidator{validator: validator.New()}
	
	e.File("/", "static/index.html")
	e.File("/style", "static/style.css")
	e.File("/script", "static/script.js")
	e.GET("/tasks", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Tasks)
	})
	e.POST("/task", func(c echo.Context) error {
		task := TaskAddRequest{}
		if err := c.Bind(&task); err != nil {
			e.Logger.Info(err)
			return c.NoContent(http.StatusBadRequest)
		}
		if err := c.Validate(task); err != nil {
			e.Logger.Info(err)
			return c.NoContent(http.StatusBadRequest)
		}
		
		Tasks[TaskId] = NewTask(task.Address)
		Tasks[TaskId].Start()
		TaskId += 1
		return c.NoContent(http.StatusOK)
	})
	
	e.Logger.Fatal(e.Start(":80"))
}
