env GOOS=linux GOARCH=amd64 go build -o builds/linux/git2git
env GOOS=windows GOARCH=amd64 go build -o builds/windows/git2git.exe
env GOOS=darwin GOARCH=amd64 go build -o builds/mac/git2git
