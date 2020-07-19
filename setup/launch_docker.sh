sudo docker build -t process_manager .
sudo docker run --gpus all -it --privileged -d  -p 5983:5983 -p 3316:3306 --name server process_manager:latest
sudo docker exec -it server /bin/bash
