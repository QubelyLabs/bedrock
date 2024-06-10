# Qubely

## Bedrock

### Introduction
A collection of reusable utilities for building APIs 

### Usage
```bash
# Published version
go get github.com/QubelyLabs/bedrock

# Local development
git clone github.com/QubelyLabs/bedrock
go mod edit -replace=github.com/QubelyLabs/bedrock@v0.0.0-unpublished=../bedrock
go get github.com/QubelyLabs/bedrock@v0.0.0-unpublished
```

### Structure
1. `pkg` - reusables for config, dbs, repository controller etc

### TODO
1. Add unit tests