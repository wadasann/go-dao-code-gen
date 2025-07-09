#!/bin/bash
go install github.com/go-bindata/go-bindata

go-bindata \
  -pkg tplbin \
  -prefix templates/ \
  -o tplbin/templates.go \
  -ignore .go$ \
  -ignore .swp$ \
  -nometadata \
  -nomemcopy templates/*.tpl
