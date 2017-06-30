#include <stdlib.h>
#include <string.h>

#include "v8.h"
#include "libplatform/libplatform.h"

#include "v8binding.h"

using namespace v8;

struct template_s {
    Isolate *isolate;
    Persistent<Context> context;
    char *last_error;
};

extern "C" {
#include "_cgo_export.h"

void _set_template_err(template_s *tpl, const char *err) {
    if (tpl->last_error) {
        free(tpl->last_error);
    }

    tpl->last_error = strdup(err);
}

void init_v8() {
    V8::InitializeICUDefaultLocation("./test");
    V8::InitializeExternalStartupData("./test");

    V8::InitializeICU();
    Platform *platform = platform::CreateDefaultPlatform();
    V8::InitializePlatform(platform);
    V8::Initialize();
}

template_s *new_template(char *tpl_source) {
    template_s *tpl = new template_s;

    Isolate::CreateParams create_params;
    create_params.array_buffer_allocator = v8::ArrayBuffer::Allocator::NewDefaultAllocator();
    tpl->isolate = Isolate::New(create_params);

    tpl->last_error = nullptr;

    Locker locker(tpl->isolate);
    Isolate::Scope isolate_scope(tpl->isolate);
    HandleScope handle_scope(tpl->isolate);
    Local<Context> context = Context::New(tpl->isolate);

    Context::Scope context_scope(context);

    TryCatch try_catch;

    Local<String> source = String::NewFromUtf8(tpl->isolate, tpl_source, NewStringType::kNormal).ToLocalChecked();

    free(tpl_source);

    Local<Script> script;
    if (!Script::Compile(context, source).ToLocal(&script)) {
        String::Utf8Value exception(try_catch.Exception());
        _set_template_err(tpl, *exception);
        return tpl;
    }

    Local<Value> result;
    if (!script->Run(context).ToLocal(&result)) {
        String::Utf8Value exception(try_catch.Exception());
        _set_template_err(tpl, *exception);
        return tpl;
    }

    tpl->context.Reset(tpl->isolate, context);

    return tpl;
}

char *eval_template(template_s *tpl, char *data) {
    Locker locker(tpl->isolate);
    Isolate::Scope isolate_scope(tpl->isolate);
    HandleScope handle_scope(tpl->isolate);
    Local<Context> context = tpl->context.Get(tpl->isolate);
    Context::Scope context_scope(context);

    Local<String> source = String::NewFromUtf8(tpl->isolate, data, NewStringType::kNormal).ToLocalChecked();
    free(data);
    Local<Script> script = Script::Compile(context, source).ToLocalChecked();
    Local<Value> result = script->Run(context).ToLocalChecked();

    String::Utf8Value utf8(result);

    return strdup(*utf8);
}

char *get_template_err(template_s *tpl) {
    return tpl->last_error;
}

void destroy_template(template_s *tpl) {
    tpl->isolate->Dispose();

    if (tpl->last_error)
        free(tpl->last_error);

    delete tpl;
}

}
