# ascii-art-stylize

## Description

`ascii-art-stylize` is a Go web application that keeps the original `ascii-art-web` rendering flow and upgrades the browser experience with custom styling. Users can enter text, choose a banner, submit the form, and view the rendered ASCII art in a styled, responsive interface.

The server exposes:

- `GET /` to display the main form
- `POST /ascii-art` to render the submitted text
- `GET /css/*` to serve styles
- `GET /image/*` to serve decorative static assets

The app supports the required banners:

- `standard`
- `shadow`
- `thinkertoy`

Escaped newlines such as `Hello\nWorld` are rendered as multiple lines, and the result is shown in an HTML page with preserved spacing.

This exercise focuses on making the existing website:

- more appealing, interactive, and intuitive
- more user-friendly
- better at giving visual feedback
- responsive and consistent across screen sizes

## Authors

- Lefteris Tzokas (Ltzokas)
- Emmanouil Damaskinakis (Edamaski)
- Alexandros Rigopoulos (arigopou)


## Usage: how to run

Requirements:

- Go installed

Start the server from the project root:

```bash
go run .
```

The server listens on:

```text
http://localhost:8080
```

How to use the web app:

1. Open `http://localhost:8080` in your browser.
2. Enter text in the textarea.
3. Select one of the available banners.
4. Submit the form.
5. View the rendered ASCII art in the styled result section on the same page.

Run the test suite:

```bash
go test ./...
```

## Implementation details: algorithm

The application starts the HTTP server in `main.go`, creates a `http.ServeMux`, registers routes through `internal/handlers`, and serves static files for styling and imagery.

Project structure:

- `main.go`
  - creates the mux
  - registers routes with `handlers.Register`
  - serves `css/` and `image/` static files
  - starts the server
- `internal/handlers`
  - owns the HTTP handlers
  - renders `templates/index.html` and `templates/error.html`
  - validates request methods, routes, and banner selection
- `internal/font`
  - loads and parses banner files
- `internal/render`
  - renders normalized input into ASCII art
- `templates/`
  - defines the HTML structure for the styled home and error pages
- `css/style.css`
  - provides the responsive visual design for the interface

Rendering flow:

1. Read `text` and `banner` from the submitted form.
2. Validate that the request method and banner are allowed.
3. Normalize escaped newlines such as `\n` into actual line breaks.
4. Load the selected banner file from `assets/`.
5. Parse the banner into glyphs with `internal/font`.
6. Render the input string with `internal/render`.
7. Inject the result back into the home page template inside a `<pre>` block so spacing and line breaks are preserved.
8. Present the result inside the styled UI without changing the ASCII output.

Error handling:

- `200 OK` for successful page loads and valid renders
- `400 Bad Request` for wrong methods or invalid banner/input
- `404 Not Found` for unknown routes and missing assets
- `500 Internal Server Error` for unexpected template or server failures

## Styling goals

- keep text readable regardless of chosen colors
- preserve consistent layout and spacing
- support responsive behavior on desktop and mobile
- improve the visual feedback for form, result, and error states
- keep the Go rendering logic unchanged unless required for correctness
