package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"gowiki/util"
)

type Page struct {
	Title string
	Body []byte
}

const pagesDir = "pages/"

var templates = template.Must(
	template.ParseFiles("templates/edit.html", "templates/view.html"))


func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(pagesDir + filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := pagesDir + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("Edit handler...")
	title, err := util.GetTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err == nil {
		p = &Page{Title: title, Body: p.Body}
	} else {
		p = &Page{Title: title }
	}
	renderTemplate(w, "templates/edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("View handler...")
	title, err := util.GetTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
	}
	renderTemplate(w, "templates/view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("Save handler...")
	title, err := util.GetTitle(w, r)
	if err != nil {
		return
	}
	body :=  r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, filepath.Base(tmpl) + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("Wiki...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", util.MakeHandler(viewHandler))
	http.HandleFunc("/edit/", util.MakeHandler(editHandler))
	http.HandleFunc("/save/", util.MakeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}

