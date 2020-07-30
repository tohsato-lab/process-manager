set -e

cd $1
TARGET=`find ./ | grep launch.sh | tail -n 1`
cd $(dirname ${TARGET})
bash launch.sh 2>&1 | tee history.log

# 0番目のコマンドのシグナルをキャッチ
signal=${PIPESTATUS[0]}
if [[ $signal -eq 0 ]]; then
	exit 0
else
	exit $signal
fi