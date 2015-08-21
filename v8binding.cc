#include <stdlib.h>
#include <string.h>

#include <iostream>

#include "v8.h"
#include "libplatform/libplatform.h"

#include "v8binding.h"

using namespace v8;

class ArrayBufferAllocator: public ArrayBuffer::Allocator {
public:
	virtual void* Allocate(size_t length) {
		void* data = AllocateUninitialized(length);
		return data == NULL ? data : memset(data, 0, length);
	}

	virtual void* AllocateUninitialized(size_t length) {
		return malloc(length);
	}

	virtual void Free(void* data, size_t) {
		free(data);
	}
};

struct template_s {
	Persistent<Context> context;
	char* last_error;
};

Isolate* isolate;
ArrayBufferAllocator allocator;

extern "C" {
#include "_cgo_export.h"

void init_v8() {
	V8::InitializeICU();
	Platform* platform = platform::CreateDefaultPlatform();
	V8::InitializePlatform(platform);
	V8::Initialize();

	Isolate::CreateParams create_params;
	create_params.array_buffer_allocator = &allocator;

	isolate = Isolate::New(create_params);
}

template_s* new_template(char* tpl_source) {
	template_s* tpl = new (template_s);
	tpl->last_error = nullptr;

	Isolate::Scope isolate_scope(isolate);
	HandleScope handle_scope(isolate);
	Local<Context> context = Context::New(isolate);

	Context::Scope context_scope(context);

	TryCatch try_catch;

	Local<String> source = String::NewFromUtf8(isolate, tpl_source, NewStringType::kNormal).ToLocalChecked();

	Local<Script> script;
	if (!Script::Compile(context, source).ToLocal(&script)) {
		String::Utf8Value exception(try_catch.Exception());
		tpl->last_error = strdup(*exception);
		return tpl;
	}

	Local<Value> result;
	if (!script->Run(context).ToLocal(&result)) {
		String::Utf8Value exception(try_catch.Exception());
		tpl->last_error = strdup(*exception);
		return tpl;
	}

	tpl->context.Reset(isolate, context);

	return tpl;
}

char* eval_template(template_s* tpl, char* data) {
	Isolate::Scope isolate_scope(isolate);
	HandleScope handle_scope(isolate);
	Local<Context> context = tpl->context.Get(isolate);
	Context::Scope context_scope(context);

	Local<String> source = String::NewFromUtf8(isolate, data, NewStringType::kNormal).ToLocalChecked();
	Local<Script> script = Script::Compile(context, source).ToLocalChecked();
	Local<Value> result = script->Run(context).ToLocalChecked();

	String::Utf8Value utf8(result);

	return strdup(*utf8);
}

char* get_template_err(template_s* tpl) {
	return tpl->last_error;
}

}
