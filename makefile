travis-ci: # must on top of the file
	@echo "add travisci commands here"

run:
	find -type f | egrep -i "*.go|*.yml|*.html|*.js" | entr -r go run *.go

build-linux: clean
	mkdir -p dist
	cp scarecrow.yml dist/
	CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '-extldflags "-static"' -o dist/scarecrow

clean:
	-rm -Rf dist/

find-port:
	--lsof -n -i :9000 | grep LISTEN
