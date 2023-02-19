IFS=$'\n' #空白を含むファイル名対策

cd "$1" || exit
TARGET=$2
ARGS=$3
CONDA_ENV=$4
ROOT=$(pwd)
DIR=("$(find . -type f -name "$TARGET" -print0 | xargs --null -n 1 echo)")

if [ ${#DIR[@]} -eq 0 ]; then
  echo 'target file is not found!!' | tee "$ROOT/history.log"
  exit 1
fi
if [ ${#DIR[@]} -ne 1 ]; then
  echo "${#DIR[@]} target files exist!" | tee "$ROOT/history.log"
  exit 1
fi
cd "$(dirname "${DIR[0]}")" || exit

source /opt/conda/etc/profile.d/conda.sh
conda activate "$CONDA_ENV"
if [ -z "$ARGS" ]; then
  python -u "$TARGET" >>"$ROOT/history.log" 2>&1
else
  python -u "$TARGET" "$ARGS" >>"$ROOT/history.log" 2>&1
fi

# 0番目のコマンドのシグナルをキャッチ
signal=${PIPESTATUS[0]}
if [[ $signal -eq 0 ]]; then
  exit 0
else
  exit "$signal"
fi
