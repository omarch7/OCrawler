package main

import (
	"io/ioutil"
	"log"
	"os"
)

type Logging struct{
	Trace *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
}

type LoggingType int

const(
	LogTrace   LoggingType = 1 + iota
	LogInfo
	LogWarning
	LogError
)

func NewLogging() *Logging  {
	logging := new(Logging)
	flags := log.Ldate | log.Ltime | log.Lshortfile
	logging.Trace = log.New(ioutil.Discard, "TRACE: ", flags)
	logging.Info = log.New(os.Stdout, "INFO: ", flags)
	logging.Warning = log.New(os.Stdout, "WARNING: ", flags)
	logging.Error = log.New(os.Stderr, "ERROR: ", flags)
	return logging
}

func (logging *Logging) Printf(t LoggingType, format string, v ...interface{})  {
	switch t {
	case LogTrace:
		logging.Trace.Panicf(format, v)
	case LogInfo:
		logging.Info.Printf(format, v)
	case LogWarning:
		logging.Warning.Printf(format, v)
	case LogError:
		logging.Error.Fatalf(format, v)
	}
}

func (logging *Logging) Print(t LoggingType, v ...interface{})  {
	switch t {
	case LogTrace:
		logging.Trace.Panic(v)
	case LogInfo:
		logging.Info.Print(v)
	case LogWarning:
		logging.Warning.Print(v)
	case LogError:
		logging.Error.Fatal(v)
	}
}

func (logging *Logging) Println(t LoggingType, v ...interface{})  {
	switch t {
	case LogTrace:
		logging.Trace.Panicln(v)
	case LogInfo:
		logging.Info.Println(v)
	case LogWarning:
		logging.Warning.Println(v)
	case LogError:
		logging.Error.Fatalln(v)
	}
}