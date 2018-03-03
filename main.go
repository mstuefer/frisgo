package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

var testdir string
var verbose bool

type test struct {
	Token  string `json:"x_auth_token"`
	Url    string
	Method string
}

type testResultContent struct {
	Field string
	Type  string
	Value string
}

type testResult struct {
	StatusCode     int    `json:"status_code"`
	ConnectionType string `json:"connection_type"`
	Contains       []testResultContent
}

type testDescription struct {
	Name   string
	Test   test
	Result testResult
}

func init() {
	flag.StringVar(&testdir, "tests", "./tests/", "Directory with all tests")
	flag.BoolVar(&verbose, "verbose", false, "Verbose to see why a test fails")
}

func main() {
	flag.Parse()

	fmt.Printf("\n**********\n* Frisgo *\n**********\n\n")

	succeeded := color.New(color.Faint, color.FgGreen).PrintfFunc()
	failed := color.New(color.Bold, color.FgRed).PrintfFunc()

	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		log.Fatal(err)
	}

	for n, file := range files {
		testName := strings.SplitN(file.Name(), "_", 3)
		err := execTest(file.Name())
		if err != nil {
			failed("FAIL ")
		} else {
			succeeded("  OK ")
		}
		fmt.Printf("-> #%d\t<%s>\n", n, testName[2])

		if err != nil && verbose {
			fmt.Printf("\t%s\n", err)
		}
	}

	fmt.Println()
}

func execTest(name string) error {
	var t testDescription

	file, err := ioutil.ReadFile(testdir + "/" + name)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(file)).Decode(&t); err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest(t.Test.Method, t.Test.Url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != t.Result.StatusCode {
		return fmt.Errorf("\n\tExpected StatusCode '%d', got '%d'\n", resp.StatusCode, t.Result.StatusCode)
	}

	for _, c := range t.Result.Contains {
		receivedValue := gjson.Get(string(body), c.Field)
		if strings.Compare(receivedValue.String(), c.Value) != 0 {
			return fmt.Errorf("\n\tExpected '%s' for field '%s', got '%s'\n", c.Value, c.Field, receivedValue)
		}
	}

	return nil
}
