go mod tidy
gomod2nix
cd ./build
GOOS=windows GOARCH=amd64 go build -v -o excel-automations.exe -ldflags "-s" ../main.go
cp ../default-config.yaml ./config.yaml
zip -r ./excel-automations.zip .