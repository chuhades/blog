package main

import (
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, path string, data interface{}) error {
	tpl := template.Must(template.ParseFiles("static/template/header.html", "static/template/footer.html", path))
	return tpl.ExecuteTemplate(w, "content", data)
}
