set -e

FILENAME=$1
SOURCE=$2
TARGET=$3

if [ -z $FILENAME ]; then
    echo "FILENAME is empty"
    exit 1
fi

# すでにログがある場合
echo $TARGET
if [ -d "$TARGET" ]; then
    echo "$TARGET is already exists."
    exit 0
fi

wget --mirror \
     --page-requisites \
     --span-hosts \
     --no-parent \
     --convert-links \
     --no-host-directories \
     --execute robots=off \
     --cut-dirs=1 \
     $SOURCE \
     -P .tmp/

echo "mv .tmp/$FILENAME $TARGET"
mv .tmp/$FILENAME $TARGET