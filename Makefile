# Test directory
TEST=./test

tests:
	go test ${TEST}/... --tags=${type} -v