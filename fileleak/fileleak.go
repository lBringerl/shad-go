//go:build !solution

package fileleak

import (
	"fmt"
	"os"
)

const procFilesDir = "/proc/self/fd"

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getOpenFilesSet() (map[string]struct{}, error) {
	filesSet := make(map[string]struct{})

	entries, err := os.ReadDir(procFilesDir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %w", err)
	}
	for _, entry := range entries {
		linkVal, err := os.Readlink(fmt.Sprintf("%s/%s", procFilesDir, entry.Name()))
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("os.Readlink: %w", err)
		}
		file := fmt.Sprintf("%s->%s", entry.Name(), linkVal)
		filesSet[file] = struct{}{}
	}

	return filesSet, nil
}

func VerifyNone(t testingT) {
	openFilesSetStart, err := getOpenFilesSet()
	if err != nil {
		panic(fmt.Sprintf("err: %s", err.Error()))
	}

	compareOpenedFiles := func() {
		openFilesSetEnd, err := getOpenFilesSet()
		if err != nil {
			panic(fmt.Sprintf("err: %s", err.Error()))
		}

		for file := range openFilesSetEnd {
			_, exists := openFilesSetStart[file]
			if !exists {
				t.Errorf("file leak")
			}
		}
	}

	t.Cleanup(compareOpenedFiles)
}
