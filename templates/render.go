package templates

import (
	"html/template"
	"net/http"
)

type Render struct {
	homeTemplate     *template.Template
	responseTemplate *template.Template
	errorTemplate    *template.Template
	notFoundTemplate *template.Template
}

func NewRender() *Render {
	return &Render{
		homeTemplate:     template.Must(template.ParseFiles("templates/html/home.html")),
		responseTemplate: template.Must(template.ParseFiles("templates/html/response.html")),
		errorTemplate:    template.Must(template.ParseFiles("templates/html/error.html")),
		notFoundTemplate: template.Must(template.ParseFiles("templates/html/not_found.html")),
	}
}

func (r *Render) Home(w http.ResponseWriter) {
	r.homeTemplate.Execute(w, nil)
}

func (r *Render) Response(w http.ResponseWriter, code string) {
	r.responseTemplate.Execute(w, code)
}

func (r *Render) Error(w http.ResponseWriter) {
	r.errorTemplate.Execute(w, nil)
}

func (r *Render) NotFound(w http.ResponseWriter, code string) {
	r.notFoundTemplate.Execute(w, code)
}
