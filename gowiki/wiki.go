package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Title       string
	Body        []byte
	DisplayBody template.HTML
}

var dataDir string = "data/"
var tmplDir string = "tmpl/"

func (p *Page) save() error {
	filename := dataDir + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := dataDir + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

var templates = template.Must(template.ParseFiles(tmplDir+"edit.html", tmplDir+"view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var linkRegexp = regexp.MustCompile("\\[([a-zA-Z0-9]+)\\]")

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	escapedBody := []byte(template.HTMLEscapeString(string(p.Body)))
	p.DisplayBody = template.HTML(linkRegexp.ReplaceAllFunc(escapedBody, func(str []byte) []byte {
		matched := linkRegexp.FindStringSubmatch(string(str))
		return []byte("<a href=\"/view/" + matched[1] + "\">" + matched[1] + "</a>")
	}))

	if title == "FrontPage" {
		files, err := ioutil.ReadDir("data/")
		if err != nil {
			log.Fatal(err)
		}

		p.DisplayBody += "</br></br></br>------------------------</br>"
		p.DisplayBody += "Content of this wiki : </br>"

		for _, f := range files {
			fileName := strings.Split(f.Name(), ".")[0]
			p.DisplayBody += template.HTML(
				[]byte("<a href=\"/view/" + fileName + "\">" + fileName + "</a>"))
			p.DisplayBody += "</br>"
		}
	}

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
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

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
