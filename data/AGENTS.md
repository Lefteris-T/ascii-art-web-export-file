# AGENTS.md — ascii-art-stylize (Go) Project Guidance

This document instructs any coding agent working in this repo to preserve the original `ascii-art-web` behavior while implementing the `ascii-art-stylize` exercise. The priority is to improve presentation and usability without breaking deterministic rendering, clean HTTP handling, or assignment-compliant error responses.

## 1) Project Goal

Maintain a Go HTTP server that:

- serves the main HTML page on `GET /`
- accepts form submissions on `POST /ascii-art`
- renders text as ASCII art in the browser
- serves styling and decorative static assets
- supports these banners:
  - `standard`
  - `shadow`
  - `thinkertoy`
- interprets literal `\\n` in submitted input as real newlines
- preserves rendered spacing and line breaks exactly
- presents the interface through a readable, responsive, styled UI

## 2) HTTP Contract (Must Follow)

- Supported application routes:
  - `GET /`
  - `POST /ascii-art`
  - static asset routes for CSS and images
- Required status behavior:
  - `200 OK` for valid requests
  - `400 Bad Request` for wrong methods, invalid banners, malformed requests, or invalid client input
  - `404 Not Found` for unknown routes and missing required files when mapped as not found
  - `500 Internal Server Error` for unexpected server-side failures

Must reject:

- non-`GET` requests to `/`
- non-`POST` requests to `/ascii-art`
- unsupported banner values

## 3) Template and Asset Rules

- HTML templates must live in the root `templates/` directory.
- Expected templates:
  - `templates/index.html`
  - `templates/error.html`
- CSS must live in `css/` and be linked from templates.
- Static images may live in `image/` and must be served without breaking existing routes.
- Banner files live in `assets/`:
  - `assets/standard.txt`
  - `assets/shadow.txt`
  - `assets/thinkertoy.txt`

Preserve whitespace exactly:

- do not trim rendered lines
- display output in `<pre>` or equivalent

Preserve readability:

- keep foreground/background contrast sufficient for text and ASCII output
- avoid style changes that obscure the rendered result
- keep the UI usable on narrow screens

## 4) Repo Architecture

- `main.go`
  - create the mux
  - call `handlers.Register`
  - serve `/css/` and `/image/`
  - keep server startup simple
- `internal/handlers/`
  - own the HTTP handlers for `/` and `/ascii-art`
  - keep template rendering and HTTP error responses out of `main.go`
- `internal/font/`
  - parse banner files into glyph maps
- `internal/render/`
  - render text line-by-line using parsed glyphs
- `templates/`
  - define the form page with in-page result rendering and the error page
- `css/`
  - define the visual system, layout, spacing, and responsive behavior

Guideline: keep `main` focused on startup, keep handlers in `internal/handlers`, and keep rendering logic reusable outside the HTTP layer.

## 5) Rendering Rules

- Read text from the submitted form.
- Normalize literal `\\n` to `\n` before rendering.
- Split normalized input by newline.
- For each logical line:
  - empty line => emit one newline
  - non-empty line => emit 8 ASCII-art rows
- Concatenate glyph rows left-to-right for each rune.
- Return deterministic errors for unsupported runes or invalid input.

## 6) Testing Requirements

- Use `go test ./...` as the baseline verification command.
- Keep coverage for:
  - banner parsing
  - rendering behavior
  - HTTP handlers and route behavior
- Manually verify styling behavior in the browser after template or CSS changes.
- HTTP tests should use:
  - `testing`
  - `net/http/httptest`

Cover at least:

- `GET /` returns `200`
- valid `POST /ascii-art` returns `200`
- invalid method returns the intended error status
- invalid banner returns `400`
- unknown route returns `404`
- missing template or internal failure path is handled
- response body contains rendered output
- static assets load correctly
- layout remains usable on mobile widths

## 7) Development Rules For Agents

Must do:

- preserve exact ASCII-art output behavior
- use only Go standard library packages
- keep error handling explicit and non-panicking
- keep routes, templates, and rendering responsibilities separated
- update README and project docs when project behavior changes
- treat CSS and template changes as product work, not cosmetic-only noise
- keep the UI consistent between the home page and error page

Must not do:

- do not reintroduce CLI-first behavior as the main product contract
- do not silently accept unsupported banners
- do not move rendering logic into `main.go`
- do not trim banner or rendered output
- do not sacrifice readability for visual effects
- do not add external frontend frameworks or packages

## 8) Verification Checklist

- server starts and listens on the documented port
- `/` renders the form page
- `/ascii-art` re-renders the home page with valid output for all three banners
- CSS is loaded and applied
- `Hello\nWorld` is rendered as two logical lines
- output spacing is preserved in the browser
- result and error states are clearly readable
- layout remains functional on small screens
- bad methods and bad routes return the right status
- project structure still matches the assignment
