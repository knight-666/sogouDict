package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestAll(t *testing.T) {
	a := "?\\//::*?\"<>||"
	reg := regexp.MustCompile(`[\\\.\*\?\|/:"<>]`)
	name := reg.ReplaceAllString(a, "_")

	fmt.Println(name)
}
