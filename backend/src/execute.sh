cd "$1" || exit
TARGET=$2
CONDA_ENV=$3
ROOT=$(pwd)
DIR=($(find ./ | grep "$TARGET"))

if [ ${#DIR[@]} -eq 0 ]; then
  echo 'target file is not found!!' | tee "$ROOT/history.log"
  exit 1
fi
if [ ${#DIR[@]} -ne 1 ]; then
  echo "${#DIR[@]} target files exist!" | tee "$ROOT/history.log"
  exit 1
fi
cd $(dirname "${DIR[0]}") || exit

source /opt/anaconda3/etc/profile.d/conda.sh
conda activate "$CONDA_ENV"
python -u "$TARGET" >> "$ROOT/history.log" 2>&1

# 0番目のコマンドのシグナルをキャッチ
signal=${PIPESTATUS[0]}
if [[ $signal -eq 0 ]]; then
  exit 0
else
  exit "$signal"
fi
