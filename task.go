package main

import (
	"bufio" 
	"fmt"
	"os/exec"
	"strings"
)

type TaskStatus int

const (
	TASK_STATUS_QUEUED TaskStatus = iota
	TASK_STATUS_DOWNLOADING
	TASK_STATUS_COMPLETE
	TASK_STATUS_FAIL
)

type Task struct {
	Name string `json:"name"`
	Address  string `json:"address"`
	Status   TaskStatus `json:"status"`
	Info     string `json:"info"`
	rawProgress string 
	Progress *TaskProgress `json:"progress"`
	fullLog  string
}

type TaskProgress struct {
	Total      string `json:"total"`
	Current    string `json:"current"`
	Speed      string `json:"speed"`
	Percentage string `json:"percentage"`
	TimeLeft   string `json:"time_left"`
}

func NewTask(videoAddress string) *Task {
	return &Task{Address: videoAddress}
}

func (t *Task) Start() {
	go func() {
		var (
			log     strings.Builder
			err     error
			line    string
			buf     []byte
			linenum int
		)

		cmd := exec.Command("./annie", t.Address)
		std, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			t.fullLog += fmt.Sprintf("\n%s", err)
			t.Status = TASK_STATUS_FAIL
		}

		t.Status = TASK_STATUS_DOWNLOADING
		reader := bufio.NewReader(std)
		for err == nil {
			buf, err = reader.ReadBytes(13)
			line = string(buf)

			if linenum == 0 {
				t.Info = line
			} else {
				t.rawProgress = strings.TrimLeft(line, " ")
				t.ParseProgress()
			}

			log.WriteString(line)
			linenum += 1
		}

		log.WriteString(fmt.Sprintf("\nexited with %s", err))
		t.fullLog = log.String()
		t.Status = TASK_STATUS_COMPLETE
	}()
}

func (t *Task) ParseProgress() bool {
	raw := strings.Split(t.rawProgress, " ")
	if len(raw) < 9 {
		return false
	}
	t.Progress = &TaskProgress{
		Current:    raw[0] + " " + raw[1],
		Total:      raw[3] + " " + raw[4],
		Percentage: raw[6],
		Speed:      raw[7] + " " + raw[8],
	}
	if len(raw) >= 10 {
		t.Progress.TimeLeft = raw[9]
	}
	return true
}
