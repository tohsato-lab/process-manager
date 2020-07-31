# install nvidia runtime
curl -s -L https://nvidia.github.io/nvidia-container-runtime/gpgkey | \
  sudo apt-key add -
 distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
 curl -s -L https://nvidia.github.io/nvidia-container-runtime/$distribution/nvidia-container-runtime.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-runtime.list
 sudo apt update
 sudo apt install nvidia-container-runtime
 sudo cp ../data/daemon.json /etc/docker/daemon.json
 sudo reboot
