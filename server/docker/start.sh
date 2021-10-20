mkdir -p /var/run/mysqld/ &&
chmod 777 /var/run/mysqld/ &&
service mysql start &&
/etc/init.d/ssh start &&
su docker -c "export PATH="$PATH:/opt/anaconda3/bin"; go build; ./process-manager-server"