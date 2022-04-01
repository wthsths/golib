package gl_session

import (
	"fmt"
	"sync"
	"time"

	"github.com/teris-io/shortid"
)

var shortIDGenerator *shortid.Shortid
var fallbackGenMutex sync.Mutex

// NewID creates a new ID value to represent a unique session.
func NewID() string {
	var err error
	if shortIDGenerator == nil {
		shortIDGenerator, err = shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			return generateFallbackID()
		}
	}

	generatedID, err := shortIDGenerator.Generate()
	if err != nil {
		return generateFallbackID()
	}

	return generatedID
}

func generateFallbackID() string {
	// Return unix timestamp upon error to always have a generated ID.
	fallbackGenMutex.Lock()
	defer fallbackGenMutex.Unlock()
	time.Sleep(1 * time.Millisecond)

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	t := fmt.Sprintf("%d", millis)
	return t
}
