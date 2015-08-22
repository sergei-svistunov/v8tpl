# v8tpl

## Description
Using the library you can execute JavaScript templates in Go by V8 engine.

## How to install
Download source:
```bash
go get -d github.com/sergei-svistunov/v8tpl/
```
Go to source directory:
```bash
cd $GOPATH/src/github.com/sergei-svistunov/v8tpl
```
Install depot_tools if you haven't it (https://www.chromium.org/developers/how-tos/install-depot-tools):
```bash
git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git
export PATH=`pwd`/depot_tools:"$PATH"
```
Build and install the library:
```bash
make -j4 install
```

## How to use
```go
package main

import (
	"github.com/sergei-svistunov/v8tpl"
)

// The JS template's source must contain function `template(data)` that will be called to process data
const TEMPLATE = `
	var template = function(data) {
		return "Hello " + data.name;
	}
`

func main() {
	// Init V8 engine
	v8tpl.InitV8()

	// Create template object
	tpl, err := v8tpl.NewTemplate(TEMPLATE)
	if err != nil {
		panic(err)
	}

	// Prepare data
	data := map[string]string{
		"name": "world",
	}

	// Process it
	res, err := tpl.Eval(data)
	if err != nil {
		panic(err)
	}

	println(res)
}
```
