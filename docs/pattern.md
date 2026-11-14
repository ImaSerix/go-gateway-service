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
- A path with `.` or `..` that came from a value is rejected: `400` `bad_request`, no request is sent. Escaping does not change such a value, and the service could fold `/public/../admin` into `/admin`. The rule has two parts:
  - the value itself is `.` or `..`, wherever in the path it stands (`/public/v{query:id}` with `id=..` is rejected too);
  - the path segment became `.` or `..` after rendering although the value is different: `/public/.{context:id}/info` with an empty `id`.
- A segment `.` or `..` written in the config as a whole is not checked. In the query `.` and `..` are regular values.
- The `400` goes first: if the same URL also has a template that cannot be resolved, the answer is still `400`, whatever the order of the templates is. Without a dot segment a template that cannot be resolved is answered as before: `500` `internal_error` for a route upstream, `503` `auth_unavailable` for a `policy` request.

| Template | Value | Result |
| --- | --- | --- |
| `/users/{route:id}/profile` | `42#` | `/users/42%23/profile` |
| `/public/{query:id}/info` | `42?admin=true#` | `/public/42%3Fadmin=true%23/info` |
| `/public/{query:id}` | `../admin/delete` | `/public/..%2Fadmin%2Fdelete` |
| `/search?locale={query:locale}` | `lv&admin=true` | `/search?locale=lv%26admin%3Dtrue` |
| `/public/{query:id}` | `..` | `400` `bad_request` |
| `/public/.{context:id}/info` | empty | `400` `bad_request` |
| `/public/{query:id}/{query:missing}` | `id` is `..`, `missing` is not sent | `400` `bad_request` |
| `/public/{query:missing}` | `missing` is not sent | `500` `internal_error` (`503` `auth_unavailable` in `policy`) |

Letters, digits, `-`, `_`, `.` inside a longer value and `~` are passed unchanged, so regular ids and UUIDs look the same as before.

A route value is encoded once. For a request path with its own encoding (`/users/a%2Fb`) the router gives `{route:key}` still encoded (`a%2Fb`); in a URL it is decoded before escaping, so the upstream gets `/users/a%2Fb`, not `/users/a%252Fb`, and `/users/%2e%2e` is rejected as `..`.

Transforms and `policy.store` do not build a URL from a string, nothing of the above applies to them: `query_params` encodes the value itself, headers and body fields get the value as it is. `{route:key}` in a transform is the value as the router gives it (`a%2Fb` for `/users/a%2Fb`).

## Example

```yaml
transforms:
  header:
    X-User-ID: "{route:id}"
```
