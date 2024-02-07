package logger

import (
	zl "github.com/rs/zerolog"
	zll "github.com/rs/zerolog/log"
	"io"
	"os"
)

type Logs struct {
	LogFile zl.Logger
	LogScr  zl.Logger
	Log2Way zl.Logger
	file    *os.File
}

func InitLogger() *Logs {

	logfile := createLogFile()

	zl.TimeFieldFormat = "02.01.2006 15:04:05"
	logs := &Logs{
		LogFile: zl.New(logfile).With().Timestamp().Logger(),
		LogScr:  zl.New(os.Stdout).With().Timestamp().Logger(),
		Log2Way: zl.New(io.MultiWriter(logfile, os.Stdout)).With().
			Timestamp().Logger(),
	}
	return logs
}

func StopLogger(l *Logs) {

	closeLogFile(l.file)
}

func createLogFile() *os.File {

	f, err := os.OpenFile("server.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zll.Fatal().Err(err)
	}
	return f
}

func closeLogFile(f *os.File) {

	err := f.Close()
	if err != nil {
		zll.Fatal().Err(err)
	}
}
