.PHONY: generated

GENERATED_FILE=levels.gen.go

generated:
	go run gen/gen.go --levels=Trace,Verbose,Debug,Info,Warning,Error,Critical,Fatal >| ${GENERATED_FILE}
	go fmt ${GENERATED_FILE}