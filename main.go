package main

import (
	elogrus "github.com/dictor/echologrus"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	Tasks         map[int]*Task = map[int]*Task{}
	TaskId        int           = 0
	Logger        elogrus.EchoLogger
	CurrentConfig *Config
	gitHash       string = "N/A"
	gitTag        string = "N/A"
	buildDate     string = "N/A"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	Logger = elogrus.Attach(e)
	e.Validator = &CustomValidator{validator: validator.New()}
	Logger.Infof("web-annie %s (%s) : %s\n", gitTag, gitHash, buildDate)
	
	successConfig := false
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		Logger.Warnln("Cannot found config file 'config.yaml'")
	} else {
		rawYaml, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			Logger.Errorf("Error is caused while reading config : %s\n", err)
		}
		rawConfig := Config{}
		if err := yaml.Unmarshal(rawYaml, &rawConfig); err != nil {
			Logger.Errorf("Error is caused while decode config : %s\n", err)
		}
		CurrentConfig = &rawConfig
		successConfig = true
	}
	if !successConfig {
		Logger.Warnln("Using default config")
		CurrentConfig = &DefaultConfig
	}
	if _, err := os.Stat(CurrentConfig.DownloadDirectory); os.IsNotExist(err) {
		os.MkdirAll(CurrentConfig.DownloadDirectory, 0775)
	}

	e.File("/", "static/index.html")
	e.File("/style", "static/style.css")
	e.File("/script", "static/script.js")
	e.GET("/tasks", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Tasks)
	})
	e.POST("/tasks", func(c echo.Context) error {
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
	e.DELETE("/tasks/:id", func(c echo.Context) error {
		reqId := c.Param("id")
		if reqId != "complete" {
			id, err := strconv.Atoi(reqId)
			if err != nil {
				e.Logger.Info(err)
				return c.NoContent(http.StatusBadRequest)
			}

			if target, exist := Tasks[id]; !exist {
				return c.NoContent(http.StatusNotFound)
			} else {
				target.Stop()
				delete(Tasks, id)
				return c.NoContent(http.StatusOK)
			}
		} else {
			deleteCnt := 0
			for i, task := range Tasks {
				if task.Status == TASK_STATUS_COMPLETE {
					task.Stop()
					delete(Tasks, i)
					deleteCnt += 1
				}
			}
			return c.JSON(http.StatusOK, map[string]int{"count": deleteCnt})
		}
	})
	e.Logger.Fatal(e.Start(CurrentConfig.ListenAddress))
}
