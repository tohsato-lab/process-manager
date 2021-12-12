export PASSWD="docker"
SOURCE=$1
TARGET=$2
PORT=8022
sshpass -p docker rsync -aloprv -e "ssh -p $PORT -o 'StrictHostKeyChecking no'" "$SOURCE" "$TARGET"