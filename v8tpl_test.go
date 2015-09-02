package v8tpl

import (
	"sync"
	"testing"
)

const TEMPLATE = `
	var template = function(data) {
		return "Hello " + data.name;
	}
`

func TestTemplate(t *testing.T) {
	tpl, err := NewTemplate(TEMPLATE)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{
		"name": "world",
	}

	res1, err := tpl.Eval(data)
	if err != nil {
		t.Error(err)
	}

	if res1 != "Hello world" {
		t.Fatalf("Invalid returned value")
	}

	res2, err := tpl.Eval(data)
	if err != nil {
		t.Error(err)
	}

	if res1 != res2 {
		t.Errorf("Results are different")
	}

}

func BenchmarkTemplate(b *testing.B) {
	tpl, err := NewTemplate(TEMPLATE)
	if err != nil {
		b.Error(err)
	}

	data := map[string]string{
		"name": "world",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tpl.Eval(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTemplateParallel(b *testing.B) {
	tpl1, err := NewTemplate(TEMPLATE)
	if err != nil {
		b.Error(err)
	}

	tpl2, err := NewTemplate(TEMPLATE)
	if err != nil {
		b.Error(err)
	}

	data := map[string]string{
		"name": "world",
	}

	b.ResetTimer()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < b.N/2; i++ {
			_, err := tpl1.Eval(data)
			if err != nil {
				b.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < b.N/2; i++ {
			_, err := tpl2.Eval(data)
			if err != nil {
				b.Error(err)
			}
		}
	}()

	wg.Wait()
}
