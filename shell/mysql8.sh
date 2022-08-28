cd /

mkdir -p mysql/conf

mkdir -p mysql/data

mkdir -p mysql/logs

docker pull mysql:8.0

docker run -d \
-p 3305:3306 \
-e MYSQL_ROOT_PASSWORD=123456 \
--name mysql \
--restart=always \
-v /mysql/data:/var/lib/mysql \
-v /mysql/conf:/etc/mysql \
-v /mysql/mysql-files:/var/lib/mysql-files/ \
-v /mysql/logs:/logs \
mysql

docker exec -it mysql /bin/bash

mysql -uroot -p123456

use mysql

alter user 'root'@'localhost' identified by '8!yJRAJwhH6t2xaK';

FLUSH PRIVILEGES;