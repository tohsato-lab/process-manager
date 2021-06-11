set -e
for i in `ls ~/process-manager/data/programs`;
do
	if [[ -z `cat ~/db.log | grep $i` ]]; then
		sudo mv "/home/migly/process-manager/data/programs/$i" '/home/migly/process-manager/tmp'
		# echo  "~/process-manager/data/programs/$i"
	fi
done
