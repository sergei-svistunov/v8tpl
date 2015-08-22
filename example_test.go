package v8tpl_test

import (
	"github.com/sergei-svistunov/v8tpl"
)

func ExampleTemplate_Eval() {
	// The JS template's source must contain function `template(data)` that will be called by executing template
	var templateSource = `
		var template = function(data) {
			return "Hello " + data.name;
		}
	`

	// Create template object
	tpl, err := v8tpl.NewTemplate(templateSource)
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