package v8tpl

/*
#cgo CFLAGS: -I./v8/include -I./v8/
#cgo CXXFLAGS: -std=c++11 -I./v8/include -I./v8/
#cgo LDFLAGS: -Wl,--start-group ./v8/out/native/obj.target/tools/gyp/libv8_libbase.a ./v8/out/native/obj.target/tools/gyp/libv8_base.a ./v8/out/native/obj.target/tools/gyp/libv8_libplatform.a ./v8/out/native/obj.target/tools/gyp/libv8_nosnapshot.a -Wl,--end-group -lrt -ldl
#include "v8binding.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
)

func InitV8() {
	C.init_v8()
}

type Template struct {
	cTemplate *C.tpl
}

func NewTemplate(source string) (*Template, error) {
	cTpl := C.new_template(C.CString(source))
	lastError := C.get_template_err(cTpl)

	if lastError != nil {
		return nil, errors.New(C.GoString(lastError))
	}

	return &Template{
		cTemplate: cTpl,
	}, nil
}

func (t *Template) Eval(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return C.GoString(C.eval_template(t.cTemplate, C.CString("template("+string(data)+");"))), nil
}
