all:
	@echo "Usage: make OPTION"
	@echo
	@echo "The options are:"
	@echo "    protobuf - create internal test protobuf descriptor"
	@echo "    bench    - run benchmark suite"
	@echo "    profile  - write a CPU profile"
	@echo

protobuf:
	protoc --go_out=. internal/proto/test.proto

bench:
	go test -benchmem -bench .

profile:
	go test -count 5 -cpuprofile cpu.prof -bench BenchmarkProfile
