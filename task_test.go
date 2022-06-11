package main

import (
	"fmt"
	"testing"
)

func TestParseProgress(t *testing.T) {
	TestInput := []string{
		"512.00 KiB / 10.75 MiB [==>------------------------------------------------------------] 76.59 KiB p/s 4.65% 2m16s",
		"62.88 MiB / 396.86 MiB [=========>--------------------------------------------------] 134.53 KiB p/s 15.84% 42m22s",
		"Merging video parts into Nature Beautiful short video 720p HD.mp4",
	}

	for _, input := range TestInput {
		task := &Task{rawProgress: input}
		task.ParseProgress()
		fmt.Printf("%+v\n", task.Progress)
	}
}
