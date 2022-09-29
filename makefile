SHELL := /bin/bash

git-add:
	git add -A
	git commit -m "New files added"

run:
	go run main.go


# build:
# 	go build -ldflags "-X main.build=local" # setting build var to "local"