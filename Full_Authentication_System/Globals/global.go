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
	EncId string
	User Mod.User
}

var (
	Msg Message
	User Mod.User
)