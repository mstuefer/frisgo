# Frisgo

## Introduction

Automated testing of APIs.
Frisgo is written in GO, very minimalistic, with the intent of getting the job done quickly.

## Motivation

When working on a code-base which comes without tests, is not written in a testable manner, and there
is not enough time/budget to change the code.. but a bug shall still be fixed or a new feature
implemented, then it might come in handy doing some "black-box testing".

## Goal

With Frisgo you can write and run several tests towards an API implementation. The final goal is to
have n testsuits running in parallel, of which each might consists of several tests running in sequence
within the suite.

Also, it shall be possible to import different test scripts from other test tools.

## Test Filename Structure

```
TestSuiteNumber_TestNumber_Description.json
```

## Example Test

```javascript
{
  "name": "Hello World Test",
  "created_at": "date",
  "created_by": "mstuefer <mstuefer@gmail.com>",
  "modified_at": "date",
  "modified_by": "dude <mstuefer+dude@gmail.com>",
  "test": {
    "type": "web",
    "x_auth_token": "my_token",
    "url": "http://localhost:8080/helloworld",
    "method": "GET"
  },
  "result": {
    "status_code": 200,
    "conection_type": "text/json",
    "contains": [
      {
        "field": "message",
        "type": "string",
        "value": "Hello World"
      },
      {
        "field": "name",
        "type": "string",
        "value": "foo"
      },
      {
        "field": "person.lastname",
        "type": "string",
        "value": "bar"
      }
    ]
  }
}
```

## Contributors

Inspired by frisby.js

## License

The MIT License
