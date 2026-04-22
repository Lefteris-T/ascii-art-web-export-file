# PRD — ascii-art-web-export

## 1. Problem Statement

We have a working Go web application that renders ASCII art. This exercise adds export functionality so users can download the rendered ASCII art as TXT, HTML, or JSON. The feature must preserve all existing rendering and validation logic, add a new HTTP export endpoint, and ensure TXT downloads match the browser-rendered output byte-for-byte.

## 2. Users and Use Case

- Primary user: student or developer using the ASCII art web app.
- Primary goal: render text as ASCII art and download the result as a file.
- Supported export formats: TXT, HTML, and JSON.
- Secondary goal: verify that export functionality does not regress existing rendering, validation, or error handling.

## 3. Product Contract

### 3.1 Existing Routes (Preserved)

- `GET /`
  - serves the main page with render form and export control
  - returns `200 OK`

- `POST /ascii-art`
  - accepts form data with `text` and `banner`
  - validates inputs
  - renders ASCII art on same page
  - returns `200 OK` on success

### 3.2 New Export Route

- `POST /export`
  - accepts form data with `text`, `banner`, and optional `format`
  - re-renders from submitted values (matches render flow exactly)
  - returns the selected downloadable file format
  - returns `200 OK` with proper HTTP headers on success
  - returns `400 Bad Request` for invalid input

### 3.3 Supported Banners

- `standard`
- `shadow`
- `thinkertoy`

### 3.4 Input Rules (Inherited from Render)

- submitted text may contain escaped newlines such as `Hello\nWorld`
- escaped newlines must be converted to actual line breaks before rendering
- export validation rules match render validation exactly

### 3.5 Export-Specific Rules

- Export endpoint re-renders from `text` + `banner` to ensure consistency
- TXT exported content must match browser-rendered result byte-for-byte
- HTML export wraps the rendered result in a valid HTML document and escapes the result before writing it into `<pre>`
- JSON export includes the submitted text, selected banner, and rendered result
- Exported file receives correct permissions (read/write for user)
- Downloaded filenames are stable and predictable:
  - `ascii-art.txt`
  - `ascii-art.html`
  - `ascii-art.json`

## 4. Functional Requirements

### 4.1 Web Interface (Existing + Export)

- The home page must contain:
  - a text input field or textarea
  - a banner selector
  - a submit button for render
  - an export format selector
  - a download button for export
- After render submission, the same page displays the result below the form.
- The rendered result must be in a whitespace-preserving container such as `<pre>`.
- Export button/link may:
  - appear only when a result exists, or
  - always appear with backend validation handling invalid requests

### 4.2 Export Endpoint

- A new `POST /export` endpoint accepts the same core form data as `/ascii-art`:
  - `text` field: the input text
  - `banner` field: the banner type
  - `format` field: `txt`, `html`, or `json`
- Export re-renders from submitted `text` + `banner` (same as render flow)
- Missing `format` defaults to `txt`
- TXT exported content must match rendered output exactly (byte-for-byte)
- Unsupported formats return `400 Bad Request`

### 4.3 HTTP Response Headers for Export

- `Content-Type`: selected media type
  - TXT: `text/plain; charset=utf-8`
  - HTML: `text/html; charset=utf-8`
  - JSON: `application/json`
- `Content-Length`: exact size of the export body
- `Content-Disposition`: attachment with selected stable filename
- Return `200 OK` on successful export

### 4.4 Input Validation for Export

- Validate `text` field using existing validation rules
- Validate `banner` field against allowed banners
- Return `400 Bad Request` if:
  - required field is missing
  - banner is invalid or missing
  - format is unsupported
  - text is empty (follow same rule as `/ascii-art`)
  - any existing validation rule fails
- Return `500 Internal Server Error` if render fails or asset is missing

### 4.5 Banner and Rendering (Existing + Shared)

- Banner files are loaded from `assets/<banner>.txt`
- Rendering reuses internal packages:
  - `internal/font`
  - `internal/render`
- Behavior preserved:
  - exact spacing
  - exact line breaks
  - multi-line input behavior
  - printable ASCII glyph handling
- Export and render flows must produce identical output

### 4.6 Error Handling (Existing + Export)

- `200 OK`
  - valid home page request
  - valid render request
  - valid export request with correct headers
- `400 Bad Request`
  - wrong HTTP method
  - invalid banner
  - missing or invalid input
- `404 Not Found`
  - unknown route
  - missing template or asset file
- `500 Internal Server Error`
  - unexpected server-side failure
  - render failure for export request

## 5. Non-Goals

- No external packages (standard library only).
- No changes to existing rendering pipeline or font logic.
- No new banner formats.
- No REST API beyond required form submission flow.
- No client-side rendering or JavaScript logic for export.
- No compression of export content.
- No breaking changes to existing routes or behavior.

## 6. Acceptance Criteria

### 6.1 Core Export Behavior

- Export endpoint exists at `POST /export`
- Export accepts `text`, `banner`, and optional `format` fields
- Export re-renders from submitted `text` + `banner`
- TXT exported content matches rendered output byte-for-byte
- Downloaded file can be TXT, HTML, or JSON
- Exported file has correct permissions (read/write for user)

### 6.2 HTTP Response

- Successful export returns `200 OK`
- Response includes correct `Content-Type` for the selected format
- Response includes `Content-Length` with correct size
- Response includes `Content-Disposition` with the selected stable filename
- Browser downloads file with correct name and extension

### 6.3 Input Validation

- Invalid banner returns `400 Bad Request`
- Missing text returns `400 Bad Request`
- Unsupported format returns `400 Bad Request`
- Empty text follows same rule as `/ascii-art`
- Missing asset file returns `500 Internal Server Error`

### 6.4 Regression and Integration

- `GET /` still works and includes export control
- `POST /ascii-art` still works without changes
- All existing validation rules still apply
- Invalid method to export endpoint returns `400`
- Unknown routes still return `404`
- Export does not break existing error handling
- All three banners work for export

### 6.5 Code Quality

- Export logic tested (unit + integration tests)
- Export reuses existing validation and render logic (DRY)
- HTTP headers set correctly in all cases
- Code follows existing handler patterns
- Only Go standard library packages used

## 7. Implementation Approach

### 7.1 TDD Strategy

- Start with export contract and HTTP specification (tasks 1-2)
- Write failing tests before implementation (tasks 3-5)
- Implement smallest working slice (tasks 6-9)
- Add regression tests (task 13)
- Refactor after all tests pass (task 14)

### 7.2 Existing Structure (Preserved)

- `main.go`
  - creates HTTP mux
  - registers routes through `internal/handlers`
  - exposes static file handlers
  - starts `http.ListenAndServe`
- `internal/handlers` (extended with export)
  - owns `GET /`, `POST /ascii-art`, and new `POST /export`
  - validates methods, fields, and banners
- `templates/`
  - `index.html` updated with export button/link
  - `error.html` for error responses
- `internal/font`
  - load and parse banner files
- `internal/render`
  - render normalized input into ASCII art

### 7.3 New Export Layer

- `internal/export/export.go` (new)
  - Export content generation function
  - Accepts selected format, submitted text, selected banner, and final ASCII result
  - Returns export body, content type, filename, and error
  - Supports TXT, HTML, and JSON
  - Minimal and deterministic
- `internal/export/export_test.go` (new)
  - Tests for TXT byte-for-byte identity
  - Tests for HTML escaping/wrapping
  - Tests for JSON payload fields
  - Tests unsupported format behavior

### 7.4 Test Coverage (new)

- `internal/handlers/export_test.go`
  - Handler tests for export endpoint
  - Tests for HTTP headers
  - Tests for error cases
- `internal/export/export_test.go`
  - Unit tests for TXT, HTML, JSON, and unsupported format behavior

### 7.5 Rationale

- Keep export content logic separate (`internal/export/`)
- Extend handlers for HTTP behavior (reuse existing package)
- Reuse existing render, font, and validation logic
- Add export fixture data under testdata/export/
- Preserve all existing file organization
- No breaking changes to existing code

## 8. Verification Checklist

### 8.1 Automated Tests (Required)

- ✓ Export handler accepts POST with valid text and banner
- ✓ Export handler returns `200 OK` on success
- ✓ Export response includes correct `Content-Type`
- ✓ Export response includes correct `Content-Length`
- ✓ Export response includes correct `Content-Disposition`
- ✓ Export supports TXT, HTML, and JSON formats
- ✓ Export rejects unsupported formats with `400`
- ✓ Export rejects invalid banner with `400`
- ✓ Export rejects missing text with `400`
- ✓ Export rejects wrong HTTP method with `400`
- ✓ TXT export content matches rendered output exactly
- ✓ Escaped newlines in export match render flow
- ✓ All existing `/ascii-art` tests still pass
- ✓ All existing `GET /` tests still pass

### 8.2 Regression Tests

- ✓ `GET /` returns `200`
- ✓ `POST /ascii-art` with valid data returns `200`
- ✓ Invalid methods return `400`
- ✓ Invalid banners return `400`
- ✓ Unknown routes return `404`
- ✓ Missing asset returns `500`
- ✓ Error page rendering works
- ✓ Response body contains rendered ASCII art

### 8.3 Manual Verification

- ✓ Open `/` in browser
- ✓ See render form and export button/link
- ✓ Submit valid text and banner
- ✓ See result rendered on page
- ✓ Select TXT, HTML, and JSON export formats
- ✓ Click export/download
- ✓ Verify files download with `.txt`, `.html`, and `.json` extensions
- ✓ Verify downloaded TXT content matches rendered output
- ✓ Test invalid banner to export returns error
- ✓ Test missing text to export returns error
- ✓ Verify all three banners work for export

### 8.4 Exercise Requirements

- ✓ Export functionality works
- ✓ Multiple export formats (TXT, HTML, JSON)
- ✓ Proper HTTP headers (Content-Type, Content-Length, Content-Disposition)
- ✓ Website includes download button/link
- ✓ Errors handled correctly
- ✓ Code respects good practices (TDD, DRY, separation of concerns)
- ✓ Only standard Go packages used
