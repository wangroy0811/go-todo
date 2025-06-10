#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo "Error: .env file not found"
    exit 1
fi

# Function to log messages
log_message() {
    local message="$1"
    local log_type="$2"
    local timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    
    echo "[${timestamp}] $message"
    
    if [[ "$log_type" == "error" ]]; then
        echo "[${timestamp}] $message" >> logs/error.log
    fi
}

# Check MySQL connection
check_mysql() {
    log_message "Checking MySQL connection..." "info"
    
    if ! command -v mysql &> /dev/null; then
        log_message "Error: MySQL client is not installed" "error"
        return 1
    fi
    
    if mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "USE $MYSQL_DB;" &> /dev/null; then
        log_message "MySQL connection successful" "info"
        return 0
    else
        log_message "Error: Failed to connect to MySQL database" "error"
        return 1
    fi
}

# Check if a port is in use and kill the process if required
check_port() {
    local port="$PORT"
    
    log_message "Checking if port $port is in use..." "info"
    
    if lsof -i :$port -t &> /dev/null; then
        local pid=$(lsof -i :$port -t)
        log_message "Port $port is in use by process $pid. Attempting to kill..." "info"
        kill -15 $pid
        
        # Wait for the process to terminate
        sleep 2
        
        if lsof -i :$port -t &> /dev/null; then
            log_message "Failed to kill process using port $port. Trying force kill..." "info"
            kill -9 $(lsof -i :$port -t) || {
                log_message "Error: Failed to kill process using port $port" "error"
                return 1
            }
        fi
        
        log_message "Process using port $port has been terminated" "info"
    else
        log_message "Port $port is available" "info"
    fi
    
    return 0
}

# Start the application
start_app() {
    log_message "Starting application..." "info"
    
    go run main.go
    
    if [ $? -ne 0 ]; then
        log_message "Error: Application failed to start" "error"
        return 1
    fi
    
    return 0
}

# Main execution
main() {
    # Create logs directory if it doesn't exist
    mkdir -p logs
    
    # Check MySQL connection
    check_mysql || {
        log_message "Error: MySQL database is not accessible. Please check your database configuration." "error"
        exit 1
    }
    
    # Check and free port if needed
    check_port || {
        log_message "Error: Failed to free port $PORT" "error"
        exit 1
    }
    
    # Start the application
    start_app || {
        log_message "Error: Application startup failed" "error"
        exit 1
    }
}
# Run the main function
main
