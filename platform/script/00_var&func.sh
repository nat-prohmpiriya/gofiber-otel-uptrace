#!/bin/bash

# color 
YELLOW='\033[0;33m'
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# server
USERSERVER="root@89.0.142.86"

# path
PROJECTPATH=""
REMOTEPROJECTPATH="/opt/dackbok"

# function
# info [timestamp] message
info(message) {
    echo -e "${BLUE}[`date +'%Y-%m-%d %H:%M:%S'`] $message${NC}"
}

# error [timestamp] message
error(message) {
    echo -e "${RED}[`date +'%Y-%m-%d %H:%M:%S'`] $message${NC}"
    exit 1
}

# warning [timestamp] message
warning(message) {
    echo -e "${YELLOW}[`date +'%Y-%m-%d %H:%M:%S'`] $message${NC}"
}

# success [timestamp] message
success(message) {
    echo -e "${GREEN}[`date +'%Y-%m-%d %H:%M:%S'`] $message${NC}"
}

# debug [timestamp] message
debug(message) {
    echo -e "${BLUE}[`date +'%Y-%m-%d %H:%M:%S'`] $message${NC}"
}

# check status
check_status(message) {
    if [ $? -eq 0 ]; then
        success "[`date +'%Y-%m-%d %H:%M:%S'`] $message"
    else
        error "[`date +'%Y-%m-%d %H:%M:%S'`] $message"
        exit 1
    fi
}

# setup_docker
setup_docker() {
    ssh $USERSERVER "sudo apt-get update"
    ssh $USERSERVER "sudo apt-get install ca-certificates curl"
    ssh $USERSERVER "sudo install -m 0755 -d /etc/apt/keyrings"
    ssh $USERSERVER "sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc"
    ssh $USERSERVER "sudo chmod a+r /etc/apt/keyrings/docker.asc"
    
    # Add repository
    ssh $USERSERVER "echo \
      \"deb [arch=\$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      \$(. /etc/os-release && echo \"\$VERSION_CODENAME\") stable\" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null"
    ssh $USERSERVER "sudo apt-get update"
}
