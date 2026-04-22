package handlers

import (
	"ascii-art-web-export-file/internal/font"
	"ascii-art-web-export-file/internal/render"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type indexPageData struct {
	Text      string
	Banner    string
	Result    string
	HasResult bool
}

type errorPageData struct {
	Code        int
	Message     string
	Title       string
	Explanation string
	Action      string
}

// Register wires the HTTP routes used by the app into the provided mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/", home)
	mux.HandleFunc("/ascii-art", asciiArt)
	mux.HandleFunc("/export", exportHandler)
}

// renderErrorPage centralizes HTML error responses and falls back to
// http.Error if the shared error template cannot be loaded.
func renderErrorPage(w http.ResponseWriter, code int, message string) {
	tmpl, err := template.ParseFiles("../../templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title, explanation, action := errorDetails(code, message)

	data := errorPageData{
		Code:        code,
		Message:     message,
		Title:       title,
		Explanation: explanation,
		Action:      action,
	}
	w.WriteHeader(code)

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}

}

func errorDetails(code int, message string) (string, string, string) {
	switch code {
	case http.StatusBadRequest:
		return "Your request could not be processed",
			"This usually means the form submission was incomplete, the selected banner was invalid, or the request used the wrong method.",
			"Go back to the home page, check your text and banner selection, then submit the form again."
	case http.StatusNotFound:
		return "The page or resource was not found",
			"The address you opened does not exist here, or the app could not find a required file for this request.",
			"Return to the home page and try again from there. If the problem continues, check that the project files are present."
	default:
		return "The server hit an unexpected problem",
			"The application could not finish this request because of an internal error.",
			"Try again in a moment. If the error keeps happening, inspect the server logs and project files."
	}
}

func renderTemplate(w http.ResponseWriter, file string, data any) error {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

// home serves the main form page and explicitly handles bad methods and
// unknown paths so the app returns assignment-specific statuses.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderErrorPage(w, http.StatusNotFound, "Not Found")
		return
	}

	if r.Method != http.MethodGet {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	data := indexPageData{
		Text:      "",
		Banner:    "standard",
		Result:    "",
		HasResult: false,
	}
	if err := renderTemplate(w, "../../templates/index.html", data); err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// asciiArt processes the submitted form, validates the chosen banner, renders
// the ASCII art, and then re-renders the home page with the result embedded.
func asciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// The exercise allows users to type literal \n in the form input.
	text := strings.ReplaceAll(r.FormValue("text"), "\\n", "\n")
	banner := r.FormValue("banner")

	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	bannerPath := filepath.Join("../../assets", banner+".txt")

	// Keep banner parsing and text rendering in the dedicated internal packages
	// so the HTTP layer stays focused on request/response handling.
	f, err := font.LoadBanner(bannerPath)
	if err != nil {
		renderErrorPage(w, http.StatusNotFound, "Not Found")
		return
	}

	result, err := render.Render(text, f)
	if err != nil {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	data := indexPageData{
		Text:      text,
		Banner:    banner,
		Result:    result,
		HasResult: true,
	}
	if err := renderTemplate(w, "../../templates/index.html", data); err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	text := strings.ReplaceAll(r.FormValue("text"), "\\n", "\n")
	banner := r.FormValue("banner")

	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	bannerPath := filepath.Join("../../assets", banner+".txt")

	f, err := font.LoadBanner(bannerPath)
	if err != nil {
		renderErrorPage(w, http.StatusNotFound, "Not Found")
		return
	}

	result, err := render.Render(text, f)
	if err != nil {
		renderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	exported := result

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)
	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(exported))))
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(exported))
	if err != nil {
		return
	}
}
