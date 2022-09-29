SHELL := /bin/bash

run:
	go run main.go

git-add:
	git add -A
	git commit -m "New files added"
