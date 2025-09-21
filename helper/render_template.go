package helper

import (
	"html/template"
	"net/http"
)

func RenderSingleTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplParsed, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplParsed.Execute(w, data)
	if err != nil {
		println("Error executing template:", err.Error())
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	files := []string{
		"template/layout/main.html",
		"template/static/nav/navbar.html",
		"template/static/footer/footer.html",
		tmpl,
	}

	tmplParsed, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplParsed.ExecuteTemplate(w, "main", data)
	if err != nil {
		println("Error executing template:", err.Error())
	}
}
