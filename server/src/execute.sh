set -e

cd $1
TARGET=`find ./ | grep launch.sh | tail -n 1`
cd $(dirname ${TARGET})
bash launch.sh 2>&1 | tee history.log
