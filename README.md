# url-shortener

A minimal URL shortening service written in Go using only the standard library.
It stores URLs in memory, hands out a random 6-character code for each one, and
redirects visitors from the code back to the original URL.

## Requirements

- Go 1.26.5 or newer (see `go.mod`)
- No third-party dependencies

## Running

```bash
git clone https://github.com/elnurAndroidDev/url-shortener.git
cd url-shortener
go run .
```

```
Server is running on port :8080
```

The port is hard-coded to `:8080`; there are no flags or environment variables.

## API

### `POST /shorten`

Creates a short code for a URL. The handler reads the value with
`r.FormValue("url")`, so the URL must be sent as a **form field or a query
string parameter — a JSON body will not work.**

```bash
curl -i --data-urlencode "url=https://example.com/a?b=1&c=2" \
  http://localhost:8080/shorten
```

```
HTTP/1.1 201 Created
Content-Type: text/plain; charset=utf-8
Content-Length: 24

http://shorty.url/Xveel7
```

The response is plain text. **Note that the `http://shorty.url/` prefix is
hard-coded and does not serve this application** — to test a code, request it
against your own server instead: `http://localhost:8080/Xveel7`.

| Status | Condition |
| --- | --- |
| `201 Created` | Short URL created; body holds the short URL |
| `400 Bad Request` | `missing 'url' parameter` — the `url` field was empty or absent |
| `405 Method Not Allowed` | The request was not a `POST` |

### `GET /<code>`

Redirects to the stored URL. Every other path falls through to this handler.

```bash
curl -i http://localhost:8080/Xveel7
```

```
HTTP/1.1 301 Moved Permanently
Location: https://example.com/a?b=1&c=2
Content-Type: text/html; charset=utf-8
```

Use `curl -i` to inspect the `Location` header, or `curl -L` to follow the
redirect. Because the redirect is a `301`, browsers cache it.

| Status | Condition |
| --- | --- |
| `301 Moved Permanently` | Code found; `Location` holds the original URL |
| `404 Not Found` | `URL not found` — no such code |
| `405 Method Not Allowed` | The request was not a `GET` |

## Short codes

Codes come from 6 bytes of `crypto/rand`, encoded as URL-safe base64 and
truncated to 6 characters drawn from `[A-Za-z0-9_-]`, giving 64^6 possible
codes.

## Project structure

```
main.go              wires up the store, handlers and routes, starts the server
handlers/handler.go  HTTP handlers for shortening and redirecting
models/url.go        random short-code generation
storage/store.go     in-memory map of code -> URL, guarded by a RWMutex
```

## Limitations

- Storage is in memory only; every mapping is lost when the server restarts.
- URLs are not validated — any non-empty string is accepted and stored.
- Code collisions are not checked; `Save` overwrites an existing mapping.
- The same URL submitted twice gets two different codes.
- The `http://shorty.url/` domain in the response is hard-coded.
- Every lookup request prints the requested code to stdout, including ones that 404.
