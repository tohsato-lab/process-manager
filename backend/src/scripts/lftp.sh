USER=docker
PASS=docker

IP=$1
PID=$2

if [ -z $PID ]; then
    exit 1
fi

lftp -c "open -u $USER,$PASS -p 31 ftp://$IP; mirror /home/process-manager/log/$PID ../../log/$PID; close; echo end; quit"
