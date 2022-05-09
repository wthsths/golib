package gl_logging

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func Test_Write_Cached(t *testing.T) {
	iterations := 100000
	logSubDir := "some-service"

	logBase := `/tmp`
	sep := "/"
	nl := "\n"
	if runtime.GOOS == "windows" {
		logBase = `C:\logs`
		sep = `\`
		nl = "\r\n"
	}

	logInit(logSubDir, logBase)

	for i := 0; i < iterations; i++ {
		writeLn(fmt.Sprint(i))
		modifyTime(1)
	}

	files, err := ioutil.ReadDir(logBase + sep + logSubDir)
	if err != nil {
		t.Fatal(err)
	}

	fileNames := make([]string, 0, len(files))

	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, logBase+sep+logSubDir+sep+file.Name())
		}
	}

	parseMap := make(map[int64]string, iterations)

	for _, fileName := range fileNames {
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatal(err)
		}

		strContent := string(content)

		lines := strings.Split(strContent, nl)

		for _, line := range lines {
			if line == "" {
				continue
			}

			parsedVal, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				t.Fatal(err)
			}

			parseMap[parsedVal] = line
		}
	}

	if len(parseMap) != iterations {
		t.Fatal()
	}
}
