.PHONY: generated

GENERATED_FILE=levels.go
GENERATED_TEST_FILE=levels_test.go

generated:
	go run gen/gen.go --path=./gen/levels.json >| ${GENERATED_FILE}
	go run gen/gen.go --path=./gen/levels.json --test=true >| ${GENERATED_TEST_FILE}
	go fmt ${GENERATED_FILE}
	go fmt ${GENERATED_TEST_FILE}