go mod init shop-srvs

go mod tidy

go build main.go

nohup ./main >> nohub.log 2>&1 &