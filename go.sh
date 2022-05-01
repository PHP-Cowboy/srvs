echo start...

cd /usr/local

wget https://studygolang.com/dl/golang/go1.18.1.linux-amd64.tar.gz

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz

cat export_go.txt >> ~/.bashrc

source ~/.bashrc

go env -w GOPROXY=https://goproxy.io,direct
go env -w GO111MODULE=on

