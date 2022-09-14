package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ltime "github.com/sviatvguss/otus-go-dev-pro/hw12_13_14_15_calendar/internal/time"
)

type Logger struct {
	Type      string
	Directory string
	Level     string
	file      *os.File
}

func New(t string, directory string, level string) (*Logger, error) {
	if t == "FILE" {
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			err := os.Mkdir(directory, 0o777)
			if err != nil {
				return nil, err
			}
		}
		f, err := os.OpenFile(filepath.Join(directory, "log.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatal(err)
		}

		return &Logger{
			Type:      t,
			Directory: directory,
			Level:     level,
			file:      f,
		}, nil
	} else {
		return &Logger{
			Type:  t,
			Level: level,
		}, nil
	}
}

func (l Logger) Info(msg string) {
	l.write(fmt.Sprintf("[%s][%s] %s\n", "INFO", time.Now().Format(ltime.DateTimeFormat), msg))
}

func (l Logger) Error(msg string) {
	l.write(fmt.Sprintf("[%s][%s] %s\n", "ERROR", time.Now().Format(ltime.DateTimeFormat), msg))
}

func (l *Logger) Close() {
	if l.Type == "FILE" {
		l.file.Close()
	}
}

func (l Logger) write(msg string) {
	if l.Type == "FILE" {
		l.file.WriteString(msg)
	} else {
		fmt.Println(msg)
	}
}
