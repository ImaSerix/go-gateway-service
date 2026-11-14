# Templates

The configuration supports templates in the `{source:key}` format.

Templates are resolved only in places that explicitly use a renderer, such as upstream URLs, transforms, and `policy.store`. In other places they are treated as regular strings.

## Sources

`source` defines where the resolver should look for the value.

Request renderer sources:
- `context` - values from `r.Context()`.
- `route` - chi route parameters.
- `query` - query parameters.
- `header` - request headers.

Response renderer sources for Store:
- `header` - response headers.
- `body` - top-level JSON body fields.

## Key

`key` is a string. If the key does not exist in the selected source, rendering returns an error.

## Values in a URL

A URL is built from a template in two places: `upstream.path` of a route and `upstream.path` of a `policy` check. There a value stays a value and cannot change the URL:

- Before `?` the value is escaped as a path segment: `/`, `?`, `#`, `%` and a space become `%2F`, `%3F`, `%23`, `%25`, `%20`.
- After `?` the value is escaped as a query value: `&`, `=`, `#`, `%` become `%26`, `%3D`, `%23`, `%25`.
- Text written in the config is not escaped and not checked.
- A value `.` or `..` in the path, or a value that turns a path segment into `.` or `..` (for example `{query:a}{query:b}` with two dots), is rejected: `400` `bad_request`, no request is sent. Such a value is not changed by escaping, and the service could fold `/public/../admin` into `/admin`. In the query `.` and `..` are regular values.

| Template | Value | Result |
| --- | --- | --- |
| `/users/{route:id}/profile` | `42#` | `/users/42%23/profile` |
| `/public/{query:id}/info` | `42?admin=true#` | `/public/42%3Fadmin=true%23/info` |
| `/public/{query:id}` | `../admin/delete` | `/public/..%2Fadmin%2Fdelete` |
| `/search?locale={query:locale}` | `lv&admin=true` | `/search?locale=lv%26admin%3Dtrue` |
| `/public/{query:id}` | `..` | `400` `bad_request` |

Letters, digits, `-`, `_`, `.` inside a longer value and `~` are passed unchanged, so regular ids and UUIDs look the same as before.

`{route:key}` is the decoded value of the path segment: for `/users/a%2Fb` it is `a/b`.

Transforms and `policy.store` do not build a URL from a string and are not escaped here: `query_params` encodes the value itself, headers and body fields get the value as it is.

## Example

```yaml
transforms:
  header:
    X-User-ID: "{route:id}"
```
