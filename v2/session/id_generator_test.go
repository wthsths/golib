package gl_session

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewShortID(t *testing.T) {
	testSize := 10000

	generatedIDs := make(map[string]bool, testSize)

	for i := 0; i < testSize; i++ {
		newID := NewID()
		fmt.Println(newID)
		idExists := generatedIDs[newID]
		if idExists {
			t.Fatalf("duplicated ID generated")
		}

		generatedIDs[newID] = true
	}
}

func Test_NanoSeconds_Unix_Diff(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	strDates := []string{
		"1990-01-01 00:00:00",
		"2000-01-01 00:00:00",
		"2010-01-01 00:00:00",
		"2020-01-01 00:00:00",
		"2050-01-01 00:00:00",
		"2200-01-01 00:00:00",
		"2200-01-01 00:00:01",
		"2200-01-01 00:01:00",
		"2200-01-01 01:00:00",
	}

	for _, strDate := range strDates {
		parsed, err := time.Parse(layout, strDate)
		if err != nil {
			t.Fatalf("could not parse: '%s', error: '%s'", strDate, err.Error())
		}
		mili := parsed.UnixNano() / 1000000
		fmt.Println(mili)
	}

	parsed, _ := time.Parse(layout, "2022-01-01 00:00:00")
	diff := time.Now().UnixNano() - parsed.UnixNano()
	fmt.Println(diff)
}
