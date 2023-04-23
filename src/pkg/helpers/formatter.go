package helpers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/color"
)

type Options struct {
	Label string
}

// Print Log
func PrintLog(title string, message string, options ...Options) string {
	var labelType string

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	logName := green("[server]:")

	if options == nil {
		labelType = blue(title)
	}

	// check object struct
	for _, opt := range options {
		if opt.Label == "success" {
			labelType = blue(title)
		} else if opt.Label == "warning" {
			labelType = yellow(title)
		} else if opt.Label == "error" {
			labelType = red(title)
		} else {
			labelType = blue(title)
		}
	}

	newMessage := cyan(message)
	result := fmt.Sprintf("%s %s %s", logName, labelType, newMessage)

	return result
}

// Pretty JSON
func PrettyJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
