package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Status      TaskStatus `json:"status"`
	Info        string     `json:"info"`
	rawProgress string
	Progress    *TaskProgress `json:"progress"`
	fullLog     string
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

		cmd := exec.Command("./annie", "-o", CurrentConfig.DownloadDirectory, t.Address)
		if CurrentConfig.HttpProxy {
			cmd.Env = append(
				os.Environ(),
				"HTTP_PROXY="+CurrentConfig.HttpProxyAddress,
			)
		}
		std, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			t.fullLog += fmt.Sprintf("\n%s", err)
			t.Status = TASK_STATUS_FAIL
			return
		}

		t.Status = TASK_STATUS_DOWNLOADING
		reader := bufio.NewReader(std)
		for err == nil {
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
			linenum += 1
		}
	
		cmd.Wait()
		exitDetail := fmt.Sprintf("\ntask success? = %t, exit code = %d", cmd.ProcessState.Success(), cmd.ProcessState.ExitCode())
		log.WriteString(exitDetail)
		t.Info += exitDetail
		
		exitString := fmt.Sprintf("\nexited with %s", err)
		log.WriteString(exitString)
		t.Info += exitString
		t.fullLog = log.String()
		t.Status = TASK_STATUS_COMPLETE
	}()
}

func (t *Task) ParseInfo() {
	rawInfo := strings.Trim(t.Info, " ")
	InfoLines := strings.Split(rawInfo, "\n")
	
	for _, l := range InfoLines {
		if strings.Contains(l, "Title") {
			t.Name = strings.Trim(strings.ReplaceAll(l, "Title:", ""), " ")
		}
	}
}

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
