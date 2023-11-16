echo "start..."

cd /usr/local

if  [ -e /usr/local/go1.21.0.linux-amd64.tar.gz ]; then   #这里是判断语句，-e表示进行比较结果为真则存在

echo "文件存在"

else

wget https://studygolang.com/dl/golang/go1.21.0.linux-amd64.tar.gz

fi

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

cd - #返回上一次目录

cat export_go.txt >> ~/.bashrc

source ~/.bashrc

go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
go env -w GO111MODULE=on

echo "end"
