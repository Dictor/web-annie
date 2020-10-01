package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	elogrus "github.com/dictor/echologrus"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

var (
	// Tasks is storage of Task
	Tasks map[int]*Task = map[int]*Task{}
	// TaskID is auto-increasing counter for assign unique id to added task
	TaskID int = 0
	// Logger is global logger reference to Echo object's logger
	Logger elogrus.EchoLogger
	// CurrentConfig is global config reference
	CurrentConfig *Config
	gitHash       string = "N/A"
	gitTag        string = "N/A"
	buildDate     string = "N/A"
)

// CustomValidator is struct validator for request input
type CustomValidator struct {
	validator *validator.Validate
}

// Validate is just renamed function of struct validate method
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
	e.File("/i18n", "static/i18n.js")
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

		Tasks[TaskID] = NewTask(task.Address)
		Tasks[TaskID].Start()
		TaskID++
		return c.NoContent(http.StatusOK)
	})
	e.DELETE("/tasks/:id", func(c echo.Context) error {
		reqID := c.Param("id")
		if reqID != "complete" {
			id, err := strconv.Atoi(reqID)
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
				if task.Status == TaskStatusComplete {
					task.Stop()
					delete(Tasks, i)
					deleteCnt++
				}
			}
			return c.JSON(http.StatusOK, map[string]int{"count": deleteCnt})
		}
	})
	e.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"tag": gitTag, "date": buildDate})
	})
	e.Logger.Fatal(e.Start(CurrentConfig.ListenAddress))
}
