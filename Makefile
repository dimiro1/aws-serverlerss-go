.PHONY: build package light deploy clean logs

build:
	GOOS=linux go build -o main

package: build
	zip deployment.zip main

logs:
	sls logs -f httpproxy

# Light deploy, does not execute cloudformation
light: package
	sls deploy -f httpproxy

# Full deploy, execute cloudformation stack
deploy: package
	sls deploy

clean:
	@rm -f deployment.zip
	@rm -f main
