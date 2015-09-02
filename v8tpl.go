/*
Package v8tpl enables to evaluate JavaScript templates in Go.

JavaScript template must implement function `template(data)` that will be called. Simple template:
	var template = function(data) {
		return "Hello " + data.name;
	}
*/
package v8tpl

/*
#cgo CFLAGS: -I${SRCDIR}/v8/include -I${SRCDIR}/v8/
#cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/v8/include -I${SRCDIR}/v8/
#cgo LDFLAGS: -Wl,--start-group ${SRCDIR}/v8/out/native/obj.target/tools/gyp/libv8_libbase.a ${SRCDIR}/v8/out/native/obj.target/tools/gyp/libv8_base.a ${SRCDIR}/v8/out/native/obj.target/tools/gyp/libv8_libplatform.a ${SRCDIR}/v8/out/native/obj.target/tools/gyp/libv8_nosnapshot.a -Wl,--end-group -lrt -ldl
#include "v8binding.h"
#include "stdlib.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"runtime"
	"unsafe"
)

func init() {
	C.init_v8()
}

// Template is a class for evaluating JS templates.
type Template struct {
	cTemplate *C.tpl
}

// NewTemplate is a constructor.
func NewTemplate(source string) (*Template, error) {
	cTpl := C.new_template(C.CString(source))
	lastError := C.get_template_err(cTpl)

	if lastError != nil {
		return nil, errors.New(C.GoString(lastError))
	}

	template := &Template{
		cTemplate: cTpl,
	}

	runtime.SetFinalizer(template, func(t *Template) {
		C.destroy_template(t.cTemplate)
	})

	return template, nil
}

// Eval evaluate JS template and return string or error if something went wrong.
func (t *Template) Eval(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	cStr := C.eval_template(t.cTemplate, C.CString("template("+string(data)+");"))
	res := C.GoString(cStr)
	C.free(unsafe.Pointer(cStr))

	return res, nil
}
