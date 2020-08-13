set -e

cd $1
TARGET=$2
CONDA_ENV=$3
ROOT=`pwd`
DIR=`find ./ | grep "$TARGET" | tail -n 1`
if [[ -z $DIR ]]; then
	echo 'targetfile not found!!' | tee "$ROOT/history.log"
	exit 1
fi
cd $(dirname ${DIR})

source /opt/anaconda3/etc/profile.d/conda.sh
conda activate "$CONDA_ENV"
python "$TARGET" 2>&1 | tee "$ROOT/history.log"

# 0番目のコマンドのシグナルをキャッチ
signal=${PIPESTATUS[0]}
if [[ $signal -eq 0 ]]; then
	exit 0
else
	exit $signal
fi