set -e

sudo docker build -t process_manager .
sudo docker run -v /home/migly/LAB/data:/process-manager-server/dataset --gpus all -it --privileged -d -p 5983:5983 -p 3316:3306 --name server process_manager:latest /bin/bash
sudo docker exec -it server /bin/bash
