# devops-tool.go
go get -u github.com/spf13/cobra
go get -u github.com/shirou/gopsutil/cpu
go get -u github.com/shirou/gopsutil/mem
go get -u github.com/shirou/gopsutil/disk
go get github.com/fatih/color
go build -o cmd
sudo mv cmd /usr/local/bin/
cmd cpu
cmd memory
