package issue

import "time"

type IssueData struct {
	IssueID      string    `json:"issue_id"`
	UserID       string    `json:"user_id"`
	AuthorName   string    `json:"author_name"`
	Date         time.Time `json:"date"`
	Title        string    `json:"title" valid:"required"`
	Device       string    `json:"device" valid:"required"`
	DeviceParsed string    `json:"device_parsed"`
	Version      string    `json:"Version" valid:"required"`
	Description  string    `json:"description" valid:"required"`
	Edited       bool      `json:"Edited"`
	Status       string    `json:"status"`
	Notify       bool      `json:"allow_notify"`
	Attachment   string    `json:"attachment_url"`
}

type IssueDataPOST struct {
	Title       string `json:"title" valid:"required" example:"Random Reboot"`
	Device      string `json:"device" valid:"required" example:"mido"`
	Version     string `json:"Version" valid:"required" example:"14.3"`
	Description string `json:"description" valid:"required" example:"Random reboot during use time"`
}
