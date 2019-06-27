package mimic

import (
	"fmt"
	"os"
	"testing"
)

const directoryName = "go-mimic"

// Path to data for testing.
var dataPath string

func TestMain(m *testing.M) {
	path, err := absPathFinder(directoryName, "data")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataPath = path
	os.Exit(m.Run())
}
