#!/bin/bash

# Agent service installation script
# Function: Check if agent service is running normally, if not, upload files and install to systemd

set -e

# Configuration variables
SERVICE_NAME="lai-agent"
AGENT_BINARY_NAME="agent"
INSTALL_DIR="/opt/lai-panel"
BIN_DIR="${INSTALL_DIR}/bin"
SYSTEMD_DIR="/etc/systemd/system"
SERVICE_FILE="${SYSTEMD_DIR}/${SERVICE_NAME}.service"
CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${CURRENT_DIR}/.." && pwd)"
BINARY_SOURCE="${PROJECT_ROOT}/bin/${AGENT_BINARY_NAME}"

# Default parameters
MASTER_HOST="${MASTER_HOST:-localhost}"
MASTER_PORT="${MASTER_PORT:-8080}"
AGENT_NAME="${AGENT_NAME:-}"
AGENT_ADDRESS="${AGENT_ADDRESS:-}"
AGENT_PORT="${AGENT_PORT:-8081}"
BINARY_SOURCE_PATH="${BINARY_SOURCE_PATH:-}"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root user
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "Please run this script with root privileges"
        exit 1
    fi
}

# Check systemd service status
check_service_status() {
    if systemctl is-active --quiet "${SERVICE_NAME}" 2>/dev/null; then
        log_info "Service ${SERVICE_NAME} is running"
        return 0
    else
        log_warn "Service ${SERVICE_NAME} is not running or does not exist"
        return 1
    fi
}

# Check if service file exists
check_service_file() {
    if [ -f "${SERVICE_FILE}" ]; then
        log_info "Service file already exists: ${SERVICE_FILE}"
        return 0
    else
        log_warn "Service file does not exist: ${SERVICE_FILE}"
        return 1
    fi
}

# Create necessary directories
create_directories() {
    log_info "Creating necessary directories..."
    mkdir -p "${BIN_DIR}"
    mkdir -p "${INSTALL_DIR}/data"
    mkdir -p "${INSTALL_DIR}/logs"
}

# Upload/copy binary file
install_binary() {
    log_info "Installing agent binary file..."
    
    # Use provided binary path if available, otherwise use default
    if [ -n "${BINARY_SOURCE_PATH}" ] && [ -f "${BINARY_SOURCE_PATH}" ]; then
        BINARY_SOURCE="${BINARY_SOURCE_PATH}"
    elif [ ! -f "${BINARY_SOURCE}" ]; then
        log_error "Source binary file does not exist: ${BINARY_SOURCE}"
        if [ -n "${BINARY_SOURCE_PATH}" ]; then
            log_error "Also tried: ${BINARY_SOURCE_PATH}"
        fi
        log_info "Please build agent first: make agent"
        exit 1
    fi
    
    cp "${BINARY_SOURCE}" "${BIN_DIR}/${AGENT_BINARY_NAME}"
    chmod +x "${BIN_DIR}/${AGENT_BINARY_NAME}"
    log_info "Binary file installed to: ${BIN_DIR}/${AGENT_BINARY_NAME}"
}

# Create systemd service file
create_service_file() {
    log_info "Creating systemd service file..."
    
    # Build startup command
    CMD="${BIN_DIR}/${AGENT_BINARY_NAME}"
    ARGS="-master-host=${MASTER_HOST} -master-port=${MASTER_PORT}"
    
    if [ -n "${AGENT_NAME}" ]; then
        ARGS="${ARGS} -name=${AGENT_NAME}"
    fi
    
    if [ -n "${AGENT_ADDRESS}" ]; then
        ARGS="${ARGS} -address=${AGENT_ADDRESS}"
    fi
    
    cat > "${SERVICE_FILE}" <<EOF
[Unit]
Description=LAI Panel Agent Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${INSTALL_DIR}
ExecStart=${CMD} ${ARGS}
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${SERVICE_NAME}

# Environment variables
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

[Install]
WantedBy=multi-user.target
EOF

    log_info "Service file created: ${SERVICE_FILE}"
}

# Install and start service
install_service() {
    log_info "Installing systemd service..."
    
    # Reload systemd configuration
    systemctl daemon-reload
    
    # Enable service (start on boot)
    systemctl enable "${SERVICE_NAME}"
    
    # Start service
    log_info "Starting service ${SERVICE_NAME}..."
    systemctl start "${SERVICE_NAME}"
    
    # Wait a bit for service to start
    sleep 2
    
    # Check service status
    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        log_info "Service ${SERVICE_NAME} started successfully"
        systemctl status "${SERVICE_NAME}" --no-pager -l
    else
        log_error "Service ${SERVICE_NAME} failed to start"
        systemctl status "${SERVICE_NAME}" --no-pager -l
        exit 1
    fi
}

# Main function
main() {
    log_info "Starting to check agent service status..."
    
    check_root
    
    # Check if service is running normally
    if check_service_status; then
        log_info "Service is running normally, no installation needed"
        exit 0
    fi
    
    # If service file does not exist or service is not running, proceed with installation
    log_info "Service is not running normally, starting installation..."
    
    create_directories
    install_binary
    create_service_file
    install_service
    
    log_info "Installation completed!"
    log_info "Use the following commands to manage the service:"
    log_info "  Check status: systemctl status ${SERVICE_NAME}"
    log_info "  View logs: journalctl -u ${SERVICE_NAME} -f"
    log_info "  Stop service: systemctl stop ${SERVICE_NAME}"
    log_info "  Restart service: systemctl restart ${SERVICE_NAME}"
}

# Display usage information
usage() {
    cat <<EOF
Usage: $0 [OPTIONS]

Options:
    -h, --help              Show help information
    --master-host HOST      Master host address (default: localhost)
    --master-port PORT      Master port (default: 8080)
    --name NAME            Agent name
    --address ADDRESS      Agent address
    --port PORT            Agent port (default: 8081)
    --binary-path PATH     Path to agent binary file (optional)

Environment variables:
    MASTER_HOST            Master host address
    MASTER_PORT            Master port
    AGENT_NAME             Agent name
    AGENT_ADDRESS          Agent address
    AGENT_PORT             Agent port
    BINARY_SOURCE_PATH     Path to agent binary file (optional)

Examples:
    $0 --master-host 192.168.1.100 --master-port 8080 --name agent-01
    MASTER_HOST=192.168.1.100 $0
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        --master-host)
            MASTER_HOST="$2"
            shift 2
            ;;
        --master-port)
            MASTER_PORT="$2"
            shift 2
            ;;
        --name)
            AGENT_NAME="$2"
            shift 2
            ;;
        --address)
            AGENT_ADDRESS="$2"
            shift 2
            ;;
        --port)
            AGENT_PORT="$2"
            shift 2
            ;;
        --binary-path)
            BINARY_SOURCE_PATH="$2"
            shift 2
            ;;
        *)
            log_error "Unknown parameter: $1"
            usage
            exit 1
            ;;
    esac
done

# Run main function
main

