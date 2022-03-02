package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCheckStatus(t *testing.T) {
	var tests = []struct {
		urls     []string
		expected int
	}{
		{[]string{"https://google.com"}, 2},
		{[]string{"https://facebook.com", "https://cnn.com"}, 4},
	}

	// Delete old logs before testing
	err := os.Remove("website-status.log")

	if err != nil {
		println("website-status.log does not exist")
	}

	// Run test cases
	for _, test := range tests {
		checkStatus(test.urls)
		content, err := ioutil.ReadFile("website-status.log")

		if err != nil {
			panic(err)
		}

		lines := strings.Split(string(content), "\n")

		if len(lines) != test.expected {
			t.Error("TEST FAILED: " + fmt.Sprint(test.expected) + " expected, but recieved " + fmt.Sprint(len(lines)))
		}
	}
}

func BenchmarkCheckStatus4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		println(n)
		checkStatus([]string{"https://facebook.com", "https://cnn.com"})
	}
}
