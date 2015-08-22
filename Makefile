version?=4.6.80

v8libs=v8/out/native/obj.target/tools/gyp/libv8_libbase.a v8/out/native/obj.target/tools/gyp/libv8_base.a v8/out/native/obj.target/tools/gyp/libv8_libplatform.a v8/out/native/obj.target/tools/gyp/libv8_nosnapshot.a

all: $(v8libs)

v8:
	fetch v8
	cd v8 && git checkout $(version) && gclient sync

$(v8libs): v8
	$(MAKE) -C v8/ i18nsupport=off native

test: $(v8libs)
	go test -v ./...

install: *.go *.cc *.h $(v8libs)
	go install

clean:
	go clean
	rm -rf v8/out

distclean: clean
	rm -f .gclient .gclient_entries
	rm -rf v8/

.PHONY: install test clean distclean
