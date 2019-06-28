package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// Unmarshal unmarshals the data residing in the designated path (from "dependencies") into the given target.
// Note that the target should be the address of an existing struct.
func Unmarshal(target interface{}, otherPaths ...string) error {
	path, err := absPathFinder(dirName, append([]string{"example", "dependencies"}, otherPaths...)...)
	if err != nil {
		return fmt.Errorf("Error getting absolute path from %v: %v", otherPaths, err)
	}
	dataFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error opening data file from %v: %v", otherPaths, err)
	}
	defer dataFile.Close()

	err = json.NewDecoder(dataFile).Decode(target)
	if err != nil {
		return fmt.Errorf("Error decoding json data from %v: %v", otherPaths, err)
	}
	return nil
}
