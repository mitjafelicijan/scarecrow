run:
	find -type f | egrep -i "*.go|*.yml|*.js" | entr -r go run *.go
