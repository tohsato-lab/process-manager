USER=docker
PASS=docker

IP=$1
PID=$2

lftp -c "open -u $USER,$PASS ftp://$IP; mirror /log/$PID ../../log/$PID; close; echo end; quit"
