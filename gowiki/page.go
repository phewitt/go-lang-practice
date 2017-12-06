package main

import (
	"html/template"
	"io/ioutil"
	"os"
)

type Page struct {
	Title       string
	Body        []byte
	DisplayBody template.HTML
}

var pathToData = "data/"

func (p *Page) save() error {
	os.Mkdir("data", 0777)
	filename := pathToData + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(pathToData + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
