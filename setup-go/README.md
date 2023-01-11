## setup-go JFrog Pipelines Task

This Pipelines Task installs and setup Go programming language

### Example:

```yaml
- task: jfrog/setup-go@v0.0.2
  input:
    version: "1.19.3"
    cacheIntegration: "art_int"
    cacheRepository: "pipelines_cache_local"
```
