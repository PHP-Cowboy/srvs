cd /home/web/ck_backup;
git pull

echo 'start build'
npm run build
echo 'build complete,moving file'
rm -rf /home/web/ck/dist2 #删除备份文件
cp -R -f /home/web/ck/dist /home/web/ck/dist2 #当前运行项目文件备份
rm -rf /home/web/ck/dist #删除当前项目运行文件
cp -R -f dist /home/web/ck/dist #把新打包的文件拷到项目目录
echo "finish"