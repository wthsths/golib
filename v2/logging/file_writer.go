package gl_logging

import (
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
)

var pathSeparator = "/"
var newLine = "\n"
var cachedName = ""
var cachedDir = ""

var cachedLogFilename = ""
var cachedLogFile *os.File
var fileNameSuffix = ""

var timeModMinute int

func logInit(name, dir string) error {
	cachedName = name
	cachedDir = dir
	fileNameSuffix = uuid.NewString()

	if runtime.GOOS == "windows" {
		pathSeparator = `\`
		newLine = "\r\n"
	}

	logDir := dir + pathSeparator + name
	fullPath := logDir + pathSeparator + "temp_" + fileNameSuffix

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(""))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	err = os.Remove(fullPath)
	if err != nil {
		return err
	}

	return nil
}

func writeLn(content string) {
	designatedFilename := time.Now().UTC().Add(time.Minute*time.Duration(timeModMinute)).Format("20060102") + "_" + fileNameSuffix

	logDir := cachedDir + pathSeparator + cachedName
	logPath := logDir + pathSeparator + designatedFilename

	if designatedFilename != cachedLogFilename {
		if cachedLogFile != nil {
			err := cachedLogFile.Close()
			if err != nil {
				panic(err)
			}
		}

		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			f, err = os.Create(logPath)
			if err != nil {
				panic(err)
			}
		}
		content += newLine

		_, err = f.WriteString(content)
		if err != nil {
			panic(err)
		}

		cachedLogFilename = designatedFilename
		cachedLogFile = f
	} else {
		content += newLine

		_, err := cachedLogFile.WriteString(content)
		if err != nil {
			panic(err)
		}
	}
}

func modifyTime(minutes int) {
	timeModMinute += minutes
}
