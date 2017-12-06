package main

import (
	"html/template"
	"net/http"
	"regexp"
)

var pathToTemplates = "tmpl/"
var templates = template.Must(template.ParseFiles(pathToTemplates+"edit.html", pathToTemplates+"view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view|)/([a-zA-Z0-0]+)$")
var linkRegexp = regexp.MustCompile("\\[([a-zA-Z0-9]+)\\]")

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":8080", nil)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	escapedBody := []byte(template.HTMLEscapeString(string(p.Body)))

	p.DisplayBody = template.HTML(linkRegexp.ReplaceAllFunc(escapedBody, func(str []byte) []byte {
		matched := linkRegexp.FindStringSubmatch(string(str))
		out := []byte("<a href=\"/view/" + matched[1] + "\">" + matched[1] + "</a>")
		return out
	}))
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
