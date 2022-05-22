package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// UnmarshalFromFile reads file content and unmarshal to placeHolder.
// It panic if placeHolder is not pointer.
func UnmarshalFromFile(filePath string, placeHolder interface{}) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, placeHolder)
}

// IsFileOld check if file is no longer modified than given duration.
// Non exisiting file treated as old.
func IsFileOld(filePath string, duration time.Duration) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf(">>> error: %v\n", err)
		// if errors.Is(err, os.ErrNotExist) {
		// 	return true, nil
		// }
		return true
	}

	return time.Now().After(stat.ModTime().Add(duration))
}

// CreateFile writes content to filePath.
// Note that content can be GO struct (must have json tag) or map not []byte
func CreateFile(filePath string, content interface{}) error {
	bs, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, bs, 0644)
}
