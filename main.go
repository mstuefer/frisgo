package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

var testdir string

// var verbose bool

type testDescription struct {
	Name   string
	Test   test
	Result testResult
}

type test struct {
	Token  string `json:"x_auth_token"`
	Url    string
	Method string
}

type testResult struct {
	StatusCode     int    `json:"status_code"`
	ConnectionType string `json:"connection_type"`
	Contains       []testResultContent
}

type testResultContent struct {
	Field string
	Type  string
	Value string
}

func init() {
	flag.StringVar(&testdir, "tests", "./tests/", "Directory with all tests")
	// flag.BoolVar(&verbose, "verbose", false, "Verbose to see why a test fails")
}

func main() {
	flag.Parse()

	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		log.Fatal(err)
	}

	for n, file := range files {
		testName := strings.SplitN(file.Name(), "_", 3)

		fmt.Printf(
			"%3d -- %s", n+1,
			strings.TrimSuffix(testName[2], filepath.Ext(testName[2])))

		execTest(file.Name())
	}
}

func execTest(name string) {
	t := readTestDescription(name)

	req, err := http.NewRequest(t.Test.Method, t.Test.Url, nil)
	if err != nil {
		log.Print(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	evaluateTest(resp, t)
}

func readTestDescription(name string) testDescription {
	file, err := ioutil.ReadFile(testdir + "/" + name)
	if err != nil {
		log.Fatal(err)
	}

	var t testDescription
	if err := json.NewDecoder(bytes.NewReader(file)).Decode(&t); err != nil {
		log.Fatal(err)
	}

	return t
}

func evaluateTest(resp *http.Response, t testDescription) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	succeeded := color.New(color.Faint, color.FgGreen).PrintfFunc()
	failed := color.New(color.Bold, color.FgRed).PrintfFunc()

	if resp.StatusCode != t.Result.StatusCode {
		failed(" [FAIL] Expected StatusCode '%d, got '%d'\n", resp.StatusCode, t.Result.StatusCode)
		return
	}

	for _, c := range t.Result.Contains {
		receivedValue := gjson.Get(string(body), c.Field)
		if strings.Compare(receivedValue.String(), c.Value) != 0 {
			failed(" [FAIL] Expected '%s' for field '%s', got '%s'\n", c.Value, c.Field, receivedValue)
			return
		}
	}

	succeeded(" [PASS]\n")
}
