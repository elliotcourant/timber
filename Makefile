.PHONY: generated

generated:
	go run gen/gen.go --levels=Trace,Verbose,Debug,Info,Warning,Error,Critical,Fatal >| levels.gen.go
	go fmt levels.gen.go