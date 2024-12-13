module workflow:
- add `import <name>` to go file
- run `go mod tidy` to update `go.mod` with required deps (or update go.mod directly)
- run `gomod2nix` to download deps via nix