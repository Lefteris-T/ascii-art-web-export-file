# ascii-art-web-export-file

## Description

`ascii-art-web-export-file` is a Go web application that renders text as ASCII art and lets the user export the rendered result as a `.txt` file.

The project keeps the existing ASCII art web flow:

- the user opens the home page
- the user enters text
- the user selects a banner
- the app renders the ASCII art in the browser
- after a result exists, the user can download the same output as a text file

The export feature re-renders the result on the server from the submitted `text` and `banner` values. It does not trust pre-rendered ASCII from the browser. This keeps the downloaded file aligned with the backend rendering rules.

Only Go standard library packages are used.

## Routes

The server exposes:

- `GET /`
  - displays the main form
- `POST /ascii-art`
  - renders submitted text as ASCII art and shows the result on the page
- `POST /export`
  - exports the rendered ASCII art as a downloadable `.txt` file
- `GET /css/*`
  - serves stylesheet files
- `GET /image/*`
  - serves static images

Supported banners:

- `standard`
- `shadow`
- `thinkertoy`

Escaped newlines such as `Hello\nWorld` are converted to real line breaks before rendering.

## Export Behavior

The export endpoint accepts the same core form values as the render flow:

- `text`
- `banner`

On success, `POST /export` returns the ASCII art as plain text with download headers:

- `Content-Type: text/plain; charset=utf-8`
- `Content-Length: <exact byte length>`
- `Content-Disposition: attachment; filename="ascii-art.txt"`

The exported body is the ASCII output only. It is not wrapped in HTML.

The browser handles the final saved file on the user's machine. The server sends a normal downloadable text response so the saved `.txt` file can be created with the user's usual read/write permissions.

## Error Handling

The app handles website and export errors through the shared error page flow.

Expected statuses:

- `200 OK`
  - valid home page request
  - valid ASCII render request
  - valid export request
- `400 Bad Request`
  - wrong HTTP method
  - invalid banner
  - invalid render input
- `404 Not Found`
  - unknown route
  - missing required banner asset
- `500 Internal Server Error`
  - unexpected template or server failure

## Usage

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
5. View the rendered ASCII art on the same page.
6. Use the export button to download the result as `ascii-art.txt`.

Run the test suite:

```bash
go test ./...
```

If the default Go build cache is not writable in your environment, use a writable cache path:

```bash
GOCACHE=/tmp/go-build go test ./...
```

## Project Structure

- `main.go`
  - starts the HTTP server
  - registers app routes through `internal/handlers`
  - serves static `css/` and `image/` files
- `internal/handlers`
  - owns `GET /`, `POST /ascii-art`, and `POST /export`
  - validates request methods and banner values
  - renders shared error pages
  - sets export download headers
- `internal/font`
  - loads and parses banner files from `assets/`
- `internal/render`
  - converts validated text into ASCII art using the selected banner
- `internal/export`
  - keeps export content generation minimal and deterministic
  - returns final ASCII content unchanged
- `templates/`
  - contains the home page and error page templates
- `assets/`
  - contains `standard.txt`, `shadow.txt`, and `thinkertoy.txt`
- `css/`
  - contains the page styling
- `data/`
  - contains exercise notes, task plan, PRD, and AI logs

## Rendering and Export Flow

1. Read `text` and `banner` from the submitted form.
2. Convert escaped `\n` sequences into actual newline characters.
3. Validate the selected banner.
4. Load the banner file from `assets/<banner>.txt`.
5. Parse the banner glyphs with `internal/font`.
6. Render the text with `internal/render`.
7. For `/ascii-art`, inject the result into the HTML page inside a `<pre>` block.
8. For `/export`, return the rendered result as plain text with download headers.

The exported content is expected to match the rendered ASCII output byte-for-byte.
The server exports the file as an HTTP download. Since the file is saved by the user's browser, final file permissions are handled by the user's operating system. The downloaded file is created under the current user account and is readable/writable by that user according to normal OS defaults.


## Authors

- Lefteris Tzokas (Ltzokas)
