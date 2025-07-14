FROM mysql:8.0

# 设置环境变量
ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=testdb
ENV MYSQL_USER=testuser
ENV MYSQL_PASSWORD=testpassword
ENV TZ=Asia/Shanghai

# 复制自定义配置文件
COPY ./mysql.cnf /etc/mysql/conf.d/my.cnf

# 设置默认认证插件
CMD ["--default-authentication-plugin=mysql_native_password"]

# 暴露 MySQL 默认端口
EXPOSE 3306