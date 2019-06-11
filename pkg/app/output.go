package app

import "github.com/fatih/color"

// Info returns info message during app execution
func Info(format string, a ...interface{}) {
	color.Blue(format, a)
}

// Error returns error mesage during app execution
func Error(format string, a ...interface{}) {
	color.Red(format, a)
}
