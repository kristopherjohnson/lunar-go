# "make" or "make lunar" - build lunar executable
# "make run"             - build lunar executable and run it
# "make test"            - build and run unit tests
# "make format"          - format Go source code with go fmt
# "make lint"            - run go vet for static analysis
# "make check"           - run both format and lint
# "make clean"           - delete executable and test output files

DIFF:=diff

lunar: lunar.go
	go build -o lunar lunar.go

run: lunar
	./lunar
.PHONY: run

test: check test_success test_failure test_good
.PHONY: test

test_good: lunar
	./lunar --echo <test/good_input.txt 1>good_output.txt 2>/dev/null || true
	$(DIFF) test/good_output_expected.txt good_output.txt
.PHONY: test_good

test_success: lunar
	./lunar --echo <test/success_input.txt 1>success_output.txt 2>/dev/null || true
	$(DIFF) test/success_output_expected.txt success_output.txt
.PHONY: test_success

test_failure: lunar
	./lunar --echo <test/failure_input.txt 1>failure_output.txt 2>/dev/null || true
	$(DIFF) test/failure_output_expected.txt failure_output.txt
.PHONY: test_failure

format:
	go fmt lunar.go
.PHONY: format

lint:
	go vet lunar.go
.PHONY: lint

check: format lint
.PHONY: check

clean:
	- $(RM) lunar
	- $(RM) success_output.txt
	- $(RM) failure_output.txt
	- $(RM) good_output.txt
.PHONY: clean
