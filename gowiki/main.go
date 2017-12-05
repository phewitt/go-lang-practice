package main

import "fmt"

func main() {
	pageName := "TestPage"
	p1 := &Page{Title: pageName, Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage(pageName)
	fmt.Println(string(p2.Body))
}
