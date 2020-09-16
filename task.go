package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TaskStatus is status indicator consisted with limited int consts
type TaskStatus int

// const for indicating each task's status
const (
	TaskStatusQueued TaskStatus = iota
	TaskStatusDownloading
	TaskStatusComplete
	TaskStatusFail
	TaskStatusCancel
)

// Task is model of downloading task
type Task struct {
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Status      TaskStatus `json:"status"`
	Info        string     `json:"info"`
	rawProgress string
	Progress    *TaskProgress `json:"progress"`
	fullLog     string
	context     context.Context
	cancel      context.CancelFunc
}

// TaskProgress is model of task's progress detail
type TaskProgress struct {
	Total      string `json:"total"`
	Current    string `json:"current"`
	Speed      string `json:"speed"`
	Percentage string `json:"percentage"`
	TimeLeft   string `json:"time_left"`
}

// NewTask is initiator of Task struct for context assigning
func NewTask(videoAddress string) *Task {
	ctx, can := context.WithCancel(context.Background())
	return &Task{Address: videoAddress, context: ctx, cancel: can}
}

// Start starts downloading
func (t *Task) Start() {
	go func() {
		var (
			log     strings.Builder
			err     error
			line    string
			buf     []byte
			linenum int
		)

		defer func() {
			t.fullLog = log.String()
		}()

		cmd := exec.Command("./annie", "-o", CurrentConfig.DownloadDirectory, t.Address)
		if CurrentConfig.HTTPProxy {
			cmd.Env = append(
				os.Environ(),
				"HTTP_PROXY="+CurrentConfig.HTTPProxyAddress,
			)
		}
		std, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			errString := fmt.Sprintf("\n%s", err)
			t.Info += errString
			t.fullLog += errString
			t.Status = TaskStatusFail
			return
		}

		t.Status = TaskStatusDownloading
		reader := bufio.NewReader(std)
		for err == nil {
			select {
			case <-t.context.Done():
				log.WriteString("\nexit because context canceled")
				t.Status = TaskStatusCancel
				return
			default:
			}

			buf, err = reader.ReadBytes(13)
			line = string(buf)

			if linenum == 0 {
				t.Info = line
				t.ParseInfo()
			} else {
				t.rawProgress = line
				t.ParseProgress()
			}

			log.WriteString(line)
			linenum++
		}

		cmd.Wait()
		exitDetail := fmt.Sprintf("\ntask success? = %t, exit code = %d", cmd.ProcessState.Success(), cmd.ProcessState.ExitCode())
		log.WriteString(exitDetail)
		t.Info += exitDetail

		exitString := fmt.Sprintf("\nexited with %s", err)
		log.WriteString(exitString)
		t.Info += exitString

		if !CurrentConfig.IgnoreExitError && !cmd.ProcessState.Success() {
			t.Status = TaskStatusFail
		} else {
			t.Status = TaskStatusComplete
		}
	}()
}

// Stop stops (cancels) downloading
func (t *Task) Stop() {
	if t.cancel != nil {
		t.cancel()
	}
}

// ParseInfo tries parsing downloading information from annie's stdout
func (t *Task) ParseInfo() {
	rawInfo := strings.Trim(t.Info, " ")
	InfoLines := strings.Split(rawInfo, "\n")

	for _, l := range InfoLines {
		if strings.Contains(l, "Title") {
			t.Name = strings.Trim(strings.ReplaceAll(l, "Title:", ""), " ")
		}
	}
}

// ParseProgress tries parsing downloading progress from annie's stdout
func (t *Task) ParseProgress() {
	//pre processing
	rawProgress := strings.Trim(t.rawProgress, " ")
	re := regexp.MustCompile("[-=>]")
	rawProgress = string(re.ReplaceAll([]byte(rawProgress), []byte("")))
	raw := strings.Split(rawProgress, " ")
	tp := TaskProgress{}

	var (
		slashAppear bool           = false
		digitReg    *regexp.Regexp = regexp.MustCompile("[0-9]+")
	)
	for i := 0; i < len(raw); i++ {
		if strings.Contains(raw[i], "/s") {
			if len(digitReg.FindAllString(raw[i-1], -1)) > 0 {
				tp.Speed = raw[i-1] + " " + raw[i]
			}
		} else if strings.ContainsAny(raw[i], "hms") {
			if len(digitReg.FindAllString(raw[i], -1)) > 0 && !strings.Contains(raw[i], "\n") && len(raw[i]) < 12 {
				tp.TimeLeft = raw[i]
			}
		} else if strings.Contains(raw[i], "B") {
			if slashAppear {
				if len(digitReg.FindAllString(raw[i-1], -1)) > 0 {
					tp.Total = raw[i-1] + " " + raw[i]
				}
			} else {
				if len(digitReg.FindAllString(raw[i-1], -1)) > 0 {
					slashAppear = true
					tp.Current = raw[i-1] + " " + raw[i]
				}
			}
		} else if strings.Contains(raw[i], "%") {
			if len(digitReg.FindAllString(raw[i], -1)) > 0 {
				tp.Percentage = raw[i]
			}
		}
	}
	t.Progress = &tp
}
