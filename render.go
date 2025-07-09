package main

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/kenshaw/snaker"
	"github.com/moon072/go-dao-code-gen/tplbin"
)

type Indexes map[string][]string // index name -> columns

type RenderData struct {
	Pkg                  string
	Table                string
	TableLowerCamelIdent string
	TableUpperCamelIdent string
	Primary              string
	Attrs                []*AttrEntity
	UniqueIndexes        Indexes
	ShadowTables         map[string]string
	TimeFields           TimeFields
	Imports              []string
}

func renderTable(name string, data *RenderData) (content []byte, err error) {
	tpl, err := tplbin.Asset("table.tpl")
	if err != nil {
		return content, err
	}
	t, err := template.New(name).Funcs(template.FuncMap{
		"ToUpperCamel": snaker.SnakeToCamelIdentifier,
	}).Parse(string(tpl))
	if err != nil {
		return content, err
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, data)
	if err != nil {
		return
	}
	content, err = format.Source(buf.Bytes())
	return
}

func renderTableConds(name string, data *RenderData) (content []byte, err error) {
	tpl, err := tplbin.Asset("conds.tpl")
	if err != nil {
		return content, err
	}
	t, err := template.New(name).Parse(string(tpl))
	if err != nil {
		return content, err
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, data)
	if err != nil {
		return
	}
	content, err = format.Source(buf.Bytes())
	return
}

func renderInitDao(data *RenderData) (content []byte, err error) {
	tpl, err := tplbin.Asset("dao.tpl")
	if err != nil {
		return content, err
	}
	t, err := template.New("dao").Parse(string(tpl))
	if err != nil {
		return content, err
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, data)
	if err != nil {
		return
	}
	content, err = format.Source(buf.Bytes())
	return
}
