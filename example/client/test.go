package client

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nat-brown/go-mimic"
)

const dirName = "go-mimic"

func absPathFinder(dirName string, otherPaths ...string) (string, error) {
	fp, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error with with getting working directory for absolute path: %v", err)
	}
	if !strings.Contains(fp, dirName) {
		return "", fmt.Errorf("filepath %s does not contain '%s'; cannot get path to test data", fp, dirName)
	}
	split := strings.SplitAfterN(fp, dirName, 2)
	return filepath.Join(append(split[:1], otherPaths...)...), nil
}

// NewTestClient returns a new client intended for testing.
func NewTestClient(paths ...string) *http.Client {
	c := newInnerClient()
	path, err := absPathFinder(dirName, append([]string{"example"}, paths...)...)
	if err != nil {
		panic(err)
	}
	mapper, err := mimic.NewMapper(path)
	if err != nil {
		panic(err)
	}
	c.Transport = &mimic.Transport{
		Mapper: &mapper,
	}
	return c
}
