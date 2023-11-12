# devops-tool.go
go get github.com/fatih/color
go build -o cmd
sudo mv cmd /usr/local/bin/
cmd cpu
cmd memory
