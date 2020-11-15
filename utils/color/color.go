package color

import "runtime"

var Reset   = "\033[0m"
var Red     = "\033[31m"
var Green   = "\033[32m"
var Yellow  = "\033[33m"
var Blue    = "\033[34m"
var Purple  = "\033[35m"
var Gray    = "\033[37m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Gray = ""
	}
}