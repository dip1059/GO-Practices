package Globals

import (
	Mod "PracticeGoland/Models"
	"html/template"
)

type Message struct {
	Success template.HTML
	Fail template.HTML
}

type EmailGenerals struct {
	From, To, Subject, HtmlString string
}

type UserDataForEmail struct {
	EncEmail string
	User Mod.User
}

var (
	Msg Message
	User Mod.User
)

const (
	DbName string = "go_crud"
)