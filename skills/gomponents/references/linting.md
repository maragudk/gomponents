# Linting

`staticcheck` flags dot imports by default (check ST1001), and golangci-lint runs it. Whitelist the gomponents packages in `.golangci.yml`:

```yaml
version: "2"
linters:
  settings:
    staticcheck:
      dot-import-whitelist:
        - "maragu.dev/gomponents"
        - "maragu.dev/gomponents/components"
        - "maragu.dev/gomponents/html"
```
