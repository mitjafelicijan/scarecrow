travis-ci: # must on top of the file
	@echo "add travisci commands here"

run:
	find -type f | egrep -i "*.go|*.yml|*.js" | entr -r go run *.go
