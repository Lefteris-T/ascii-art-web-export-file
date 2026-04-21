# PRD — ascii-art-stylize

## 1. Problem Statement

We already have a working Go web application that renders ASCII art in the browser. This exercise is to stylize that existing website without breaking the core rendering behavior. The product must remain assignment-compliant while becoming more appealing, responsive, consistent, and easier to use.

## 2. Users and Use Case

- Primary user: student or reviewer opening the project in a browser.
- Primary goal: render text as ASCII art with `standard`, `shadow`, or `thinkertoy` through a polished web UI.
- Secondary goal: verify that styling improvements do not regress input validation, error handling, readability, or route behavior.

## 3. Product Contract

The application must preserve these routes:

- `GET /`
  - serves the main page
  - returns `200 OK`

- `POST /ascii-art`
  - accepts form data
  - validates `text` and `banner`
  - renders ASCII art
  - returns `200 OK` on success

- `GET /css/*`
  - serves CSS assets used by the interface

- `GET /image/*`
  - serves static decorative assets used by the interface

Supported banners:

- `standard`
- `shadow`
- `thinkertoy`

Input rules:

- submitted text may contain escaped newlines such as `Hello\nWorld`
- escaped newlines must be converted to actual line breaks before rendering

Presentation rules:

- the interface must use CSS
- the website must remain readable regardless of color choices
- the layout must be responsive and consistent
- styling must not alter the generated ASCII output

## 4. Functional Requirements

### 4.1 Web Interface

- The home page must contain:
  - a text input field or textarea
  - a banner selector
  - a submit button
- After a successful submission, the same page should display the rendered result below the form.
- The rendered result must be displayed in a whitespace-preserving container such as `<pre>`.
- The interface must provide a clear visual distinction between:
  - the input section
  - the action controls
  - the result section
  - the error page
- The layout must remain usable on smaller screens.

### 4.2 Request Validation

- `GET /` must reject non-`GET` methods with `400 Bad Request`.
- `POST /ascii-art` must reject non-`POST` methods with `400 Bad Request`.
- Any unsupported banner value must return `400 Bad Request`.
- Unknown routes must return `404 Not Found`.

### 4.3 Banner and Rendering

- Banner files are loaded from `assets/<banner>.txt`.
- Rendering must continue to use reusable internal packages:
  - `internal/font`
  - `internal/render`
- Rendering behavior must preserve:
  - exact spacing
  - exact line breaks
  - multi-line input behavior
  - printable ASCII glyph handling

### 4.4 Styling and UX

- CSS must be linked from the HTML templates.
- The interface should feel intentionally styled rather than default browser output.
- Form controls should communicate focus and interaction clearly.
- The result block should be easy to distinguish and easy to read.
- Error pages should provide immediate feedback and a clear path back to the form.

### 4.5 Error Handling

- `200 OK`
  - valid home page request
  - valid render request
- `400 Bad Request`
  - wrong method
  - invalid banner
  - malformed or invalid client input
- `404 Not Found`
  - unknown route
  - missing template or banner file when mapped as not found
- `500 Internal Server Error`
  - unexpected server-side failure
  - template execution/parsing failure not caused by the client

## 5. Non-Goals

- No external packages.
- No rebuild of the server-side rendering pipeline.
- No REST API beyond the required HTML form flow.
- No additional banner types.
- No client-side rendering framework.
- No changes that sacrifice readability for visual effects.

## 6. Acceptance Criteria

- The app starts as a Go HTTP server.
- Visiting `/` shows the form in a styled interface.
- Submitting valid text and a valid banner renders ASCII art successfully.
- All three banners work.
- Escaped newline input works correctly.
- Output spacing is preserved in the browser.
- CSS is applied successfully through the static route.
- The interface is visually improved and remains readable.
- The layout adapts cleanly on desktop and mobile widths.
- Invalid methods return `400`.
- Invalid banner values return `400`.
- Unknown routes return `404`.
- Internal failures are handled with `500`.
- Templates are stored in the root `templates/` directory.
- The project uses only Go standard library packages.

## 7. Implementation Approach

- `main.go`
  - create the HTTP mux
  - register routes through `internal/handlers`
  - expose static file handlers for CSS and images
  - start `http.ListenAndServe`
- `internal/handlers`
  - own the `GET /` and `POST /ascii-art` handlers
  - keep HTTP validation, template rendering, and error rendering out of `main.go`
- `templates/`
  - `index.html` for the styled form and in-page result display
  - `error.html` for styled error responses
- `css/`
  - `style.css` for responsive layout, typography, colors, and interaction states
- `internal/font`
  - load and parse banner files
- `internal/render`
  - render normalized input into ASCII art

Rationale:

- keep `main.go` focused on application startup
- keep HTTP concerns grouped in one package
- preserve reusable rendering logic
- keep template and route responsibilities separate
- confine presentation work to templates, CSS, and static assets when possible

## 8. Verification Checklist

- `GET /` returns `200`.
- `POST /ascii-art` with valid form data returns `200`.
- CSS loads from the static route.
- Invalid methods return `400`.
- Invalid banners return `400`.
- Unknown routes return `404`.
- Missing template or internal failure path returns `404` or `500` as designed.
- Response body contains rendered ASCII art.
- The result remains readable after styling changes.
- The page remains usable at mobile widths.
- Manual audit cases from the exercise render correctly.
