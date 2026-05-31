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

## Example

```yaml
transforms:
  header:
    X-User-ID: "{route:id}"
```
