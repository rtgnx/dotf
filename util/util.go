package util

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func HasScheme(s string) bool {
	return len(strings.Split(s, ":")) > 1
}

func AbsOrScheme(s string) string {
	if !HasScheme(s) {
		if !path.IsAbs(s) {
			s = path.Join(Must(os.Getwd()), s)
		}
		return fmt.Sprintf("file:%s", s)
	}
	return s
}

func Must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}
