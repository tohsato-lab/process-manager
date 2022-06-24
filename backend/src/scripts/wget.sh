set -e

FILENAME=$1
SOURCE=$2
TARGET=$3

# すでにログがある場合
if [ -d $TARGET ]; then
    echo "$FILENAME is already exists."
    exit 0
fi

if [ -d .tmp/$FILENAME ]; then
    rm -r .tmp/$FILENAME
fi
wget -r --level=0 -q --show-progress -np -nH --cut-dirs=1 -R index.html $SOURCE -P .tmp
mv .tmp/$FILENAME $TARGET