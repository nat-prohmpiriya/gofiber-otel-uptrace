#!/bin/bash

# Import variables and functions
source "$(dirname "$0")/00_var&func.sh"



# Allow Port 80/443
setup_firewall() {
    ssh $USERSERVER "sudo ufw allow 80"
    ssh $USERSERVER "sudo ufw allow 443"
}

# Install Docker
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

