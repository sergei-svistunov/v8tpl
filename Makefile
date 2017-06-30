version?=6.1.313

v8libs=v8/out/native/obj.target/src/libv8_libbase.a v8/out/native/obj.target/src/libv8_base.a v8/out/native/obj.target/src/libv8_libplatform.a v8/out/native/obj.target/src/libv8_nosnapshot.a

all: $(v8libs)

v8:
	fetch v8
	cd v8 && git checkout $(version) && gclient sync

$(v8libs): v8
	export CPPFLAGS="-fPIC" CCFLAGS="-fPIC" && $(MAKE) -C v8/ i18nsupport=off snapshot=off native

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
