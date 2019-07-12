.PHONY: generated

GENERATED_FILE=levels.go
GENERATED_TEST_FILE=levels_test.go

generated:
	go run gen/gen.go --levels=Trace,Verbose,Debug,Info,Warning,Error,Critical,Fatal >| ${GENERATED_FILE}
	go run gen/gen.go --test=true --levels=Trace,Verbose,Debug,Info,Warning,Error,Critical,Fatal >| ${GENERATED_TEST_FILE}
	go fmt ${GENERATED_FILE}
	go fmt ${GENERATED_TEST_FILE}