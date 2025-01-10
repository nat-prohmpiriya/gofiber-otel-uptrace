#!/bin/bash

# Server configurations
PROD_SERVER="root@89.0.142.86"
STAGING_SERVER="root@89.0.142.86"

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Function to install Netdata
install_netdata() {
    local server=$1
    echo -e "${GREEN}Installing Netdata on $server...${NC}"
    ssh $server "curl -fsSL https://get.netdata.cloud | bash"
}

# Function to create SSH tunnel
create_tunnel() {
    local server=$1
    local local_port=$2
    echo -e "${GREEN}Creating SSH tunnel to $server on port $local_port...${NC}"
    ssh -f -N -L $local_port:localhost:19999 $server
}

# Function to check if port is in use
check_port() {
    local port=$1
    lsof -i :$port > /dev/null
    return $?
}

# Function to open browser
open_browser() {
    local url=$1
    echo -e "${GREEN}Opening $url${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open "$url"  # macOS
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        xdg-open "$url"  # Linux
    fi
}

# Main script
main() {
    # Install Netdata on both servers
    echo "=== Installing Netdata ==="
    install_netdata $PROD_SERVER
    install_netdata $STAGING_SERVER

    # Create SSH tunnels
    echo "=== Creating SSH Tunnels ==="
    if ! check_port 19999; then
        create_tunnel $PROD_SERVER 19999
        echo "Production server monitoring: http://localhost:19999"
        # Wait for tunnel to establish
        sleep 2
        open_browser "http://localhost:19999"
    else
        echo "Port 19999 is already in use"
    fi

    if ! check_port 19998; then
        create_tunnel $STAGING_SERVER 19998
        echo "Staging server monitoring: http://localhost:19998"
        # Wait for tunnel to establish
        sleep 2
        open_browser "http://localhost:19998"
    else
        echo "Port 19998 is already in use"
    fi

    echo "=== Setup Complete ==="
    echo "To view monitoring:"
    echo "Production: http://localhost:19999"
    echo "Staging: http://localhost:19998"
}

# Run main function
main
