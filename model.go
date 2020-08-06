package main

type (
	TaskAddRequest struct {
		Address string `json:"address" validate:"required,url"`
	}
)