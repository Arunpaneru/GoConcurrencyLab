package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestPrintTask(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)

	go printTask("hello", &wg)

	wg.Wait()
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "hello") {
		t.Errorf("expected to find the hello but found %s", &result)
	}
}
