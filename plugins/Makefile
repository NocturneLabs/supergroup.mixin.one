all: build/course_record.so

.SECONDEXPANSION:
build/%.so: $$(wildcard %/*.go)
	go build -buildmode=plugin -o $@ $^

clean:
	rm -rf build
