# ASCII Art Web Export - Task Plan (SDD / TDD First)

## Goal

Build an export feature for the existing `ascii-art-web` project so that:

- the user generates ASCII art from the web form
- the app can export the generated ASCII art as a `.txt` download
- the downloaded file contains **exactly the same text** as the rendered result
- invalid export requests are rejected with the correct HTTP behavior
- the feature is added without breaking the current home page, rendering flow, and error handling

---

## Development Strategy

We will follow this order:

1. define the export behavior precisely
2. decide the HTTP contract before coding
3. write failing tests first
4. implement the smallest working slice
5. refactor only after tests are green

SDD/TDD rule for this project:

1. every new behavior starts from an explicit spec
2. every spec becomes one or more failing tests
3. implementation must satisfy only the current tested behavior
4. export output must match the rendered ASCII result **byte-for-byte**
5. no hidden business logic should live inside templates

The feature should be built in layers:

1. export contract
2. export endpoint tests
3. request validation
4. export content generation
5. download response headers
6. template integration
7. regression checks with existing render flow

---

## Suggested Project Structure

```text
.
├── main.go
├── main_test.go
├── internal/
│   ├── font/
│   │   └── font.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── render/
│   │   ├── render.go
│   │   └── render_test.go
│   └── export/
│       ├── export.go
│       └── export_test.go
├── templates/
│   ├── index.html
│   └── error.html
├── assets/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
└── testdata/
    ├── export/
    └── render/

    ## Implementation Tasks

Follow these tasks in order. Do not skip ahead before the current layer is green.

### 1. Define the exact export contract in the spec.

### 1. Define the exact export contract in the spec.

- Export route: `POST /export`
- Export accepts the same core inputs as the render flow:
  - `text`
  - `banner`
- Export does **not** trust a pre-rendered ASCII result from the client.
- **Decision:** export re-renders on the server from submitted `text` + `banner`.
- The exported file body must be exactly the ASCII-art result produced by the backend render rules.
- Export response body must be plain text only.
- Export must trigger file download with a `.txt` filename.
- Export must reuse the same validation boundaries as the render flow.

- Error behavior:
  - wrong HTTP method -> reject request
  - missing `text` or invalid `text` -> reject request according to render rules
  - missing `banner` or invalid `banner` -> reject request
  - missing banner asset file -> reject request
  - render failure -> reject request

- Consistency rules:
  - newline handling must match the existing render behavior exactly
  - no extra spaces may be added
  - no lines may be removed unless the shared render rules already do that
  - exported output must match render output byte-for-byte

- File policy:
  - download filename should be stable and predictable
  - content type must represent plain text
  - response should only be successful when export content is generated correctly
  - Empty `text` handling must follow the exact same behavior as `POST /ascii-art`.
- If `POST /ascii-art` returns an empty rendered result for empty input, `POST /export` must return that same empty result as a downloadable plain text response.

### 2. Freeze the expected HTTP behavior before implementation.

- Define the export route path as `POST /export`.
- Reject any method other than `POST`.
- Return `200 OK` on successful export.
- Return a client error status for invalid requests.
- Return a server error status for internal export/render failures.
- Set `Content-Type` to plain text for the downloadable response.
- Set `Content-Disposition` so the browser downloads a `.txt` file.
- Decide and keep one stable filename policy for the exported file.
- Keep empty `text` behavior exactly aligned with `POST /ascii-art`.
- Keep newline behavior exactly aligned with the existing render flow.
- Ensure the response body is the exact ASCII output, not HTML.

### Error behavior for export

- Reuse the same HTTP error boundaries already used by `/ascii-art`.
- Return `400 Bad Request` for:
  - wrong method
  - invalid banner
  - invalid render input
- Return `404 Not Found` when the required banner asset file is missing.
- Return `500 Internal Server Error` only for unexpected internal failures.
- Reuse the existing error page rendering flow for export failures.

### 3. Add failing export handler tests.

- Cover successful export request.
- Cover wrong HTTP method.
- Cover invalid banner.
- Cover render error path.
- Cover missing banner asset.
- Assert exact response status.
- Assert download headers are correct.
- Assert response body is plain text, not HTML.
- Assert response body equals expected ASCII output exactly.

### 4. Decide the source of truth for export output.

- Choose whether export will:
  - re-render from `text` + `banner`, or
  - receive the final ASCII result directly.
- Document the choice briefly in the repo.
- Keep the chosen model consistent across handler, tests, and template flow.

### 5. Add failing tests for export content generation.

- Cover exact pass-through of ASCII content.
- Cover newline preservation.
- Cover empty content behavior.
- Cover no added spaces.
- Cover no removed lines unless explicitly required.
- Assert byte-for-byte identity between generated result and exported body.

### 6. Implement the smallest export content layer.

- Accept the final ASCII text.
- Return it unchanged as export body.
- Keep this layer deterministic and minimal.
- Avoid mixing HTTP response logic into this layer.

### 7. Implement the export HTTP handler.

- Register the new export route.
- Validate method first.
- Read required form values.
- Reuse existing text/banner validation rules.
- Produce export content from the chosen source-of-truth path.
- Return early on invalid request.
- Keep the handler focused on request/response behavior.

### 8. Add proper download response headers.

- Set `Content-Type` to plain text.
- Set `Content-Length` to the size of the export body.
- Set `Content-Disposition` to force `.txt` download.
- Use a stable and predictable filename.
- Return success status only when export body is valid.

### 9. Reuse the existing render rules exactly.

- Reuse the same text normalization behavior.
- Reuse the same banner validation.
- Reuse the same banner loading rules.
- Reuse the same renderer behavior.
- Reuse the same error boundaries.

### 10. Add failing integration tests for parity with rendered output.

- Submit the same `text` and `banner` to the render flow and export flow.
- Compare the rendered result with the exported file body.
- Assert they are identical.

### 11. Update the template flow minimally.

- Add the export trigger to the UI only after backend tests are green.
- Decide whether the export button appears:
  - only when a result exists, or
  - always, with backend validation handling invalid requests.
- Keep template logic minimal.
- Avoid duplicating business logic inside HTML.

### 12. Add template/handler integration tests if needed.

- Cover presence of the export control in the page when expected.
- Cover correct form action/route.
- Cover correct HTTP method.
- Cover required hidden fields if the chosen design needs them.

### 13. Add regression tests for existing `/ascii-art` behavior.

- Confirm `GET /` still works.
- Confirm `POST /ascii-art` still works.
- Confirm invalid method handling still works.
- Confirm bad banner handling still works.
- Confirm missing asset behavior still works.
- Confirm existing error-page behavior still works.

### 14. Refactor after the suite is green.

- Remove duplication between render and export flows.
- Extract shared request parsing only if it improves clarity.
- Improve naming.
- Tighten package boundaries.
- Keep behavior unchanged.

### 15. Run the full test suite.

- Confirm render tests pass.
- Confirm handler tests pass.
- Confirm export tests pass.
- Confirm integration tests pass.
- Confirm regression tests pass.
- Treat the feature as incomplete until the full suite is green.