all:
	@echo "Usage: make OPTION"
	@echo
	@echo "The options are:"
	@echo "    protobuf - create internal test protobuf descriptor"
	@echo "    bench    - run benchmark suite"
	@echo

protobuf:
	protoc --go_out=. internal/proto/test.proto

bench:
	go test -benchmem -count 5 -bench .
