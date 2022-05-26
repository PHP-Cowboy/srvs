cd /

mkdir -p mysql/conf

mkdir -p mysql/data

mkdir -p mysql/logs

docker pull mysql:5.7

docker run -p 3306:3306 --name mysql57 -v /mysql/conf:/etc/mysql/conf.d -v /mysql/logs:/logs -v /mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7

docker exec -it mysql57 /bin/bash

mysql -uroot -p123456

GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'root' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' IDENTIFIED BY 'root' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' IDENTIFIED BY 'root' WITH GRANT OPTION;
FLUSH PRIVILEGES;

docker container update --restart=always mysql57