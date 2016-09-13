#!/bin/bash

go get github.com/go-martini/martini
go get github.com/go-sql-driver/mysql
#!go get github.com/martini-contrib/render
#!go get github.com/martini-contrib/sessions
go get github.com/gin-gonic/gin
go get github.com/gorilla/sessions
go build -o golang-webapp .
