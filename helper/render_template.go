package helper

import (
	"html/template"
	"net/http"
	"strings"
)

func RenderSingleTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmplParsed, err := template.New("").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Find and execute the first defined template that's not the filename
	for _, template := range tmplParsed.Templates() {
		if template.Name() != "" && template.Name() != "login.html" {
			err = tmplParsed.ExecuteTemplate(w, template.Name(), data)
			if err != nil {
				println("Error executing template:", err.Error())
				http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	files := []string{
		"template/layout/main.html",
		"template/static/nav/navbar.html",
		"template/static/footer/footer.html",
		tmpl,
	}

	// Create main template first
	mainTemplate := template.New("main").Funcs(funcMap)

	// Parse all files
	tmplParsed, err := mainTemplate.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "main" template
	err = tmplParsed.ExecuteTemplate(w, "main", data)
	if err != nil {
		println("Error executing template:", err.Error())
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
