#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# =============================================
# Configuration
# =============================================

# Environment file paths
BASE_ENV_FILE=".env.example"
ENV_FILE=".env"

# Air installation URL
AIR_INSTALL_URL="https://raw.githubusercontent.com/air-verse/air/master/install.sh"

# =============================================
# Helper Functions
# =============================================

# Get system information
get_system_info() {
  echo "🔍 Detecting system information..."

  # Get Docker version
  if command -v docker &>/dev/null; then
    DOCKER_VERSION=$(docker version --format '{{.Server.Version}}' 2>/dev/null || echo "Not installed")
  else
    DOCKER_VERSION="Not installed"
  fi

  # Detect OS type and version
  if [[ "$OSTYPE" == "darwin"* ]]; then
    OS_TYPE="macos"
    OS_VERSION=$(sw_vers -productVersion)
  else
    OS_TYPE=$(grep -w "ID" /etc/os-release 2>/dev/null | cut -d "=" -f 2 | tr -d '"' || echo "unknown")
    OS_VERSION=$(grep -w "VERSION_ID" /etc/os-release 2>/dev/null | cut -d "=" -f 2 | tr -d '"' || echo "unknown")

    # Normalize OS type
    if [[ "$OS_TYPE" =~ ^(ubuntu|alpine)$ ]]; then
      OS_TYPE="linux"
    fi
  fi
}

# Install Air for live reload
install_air() {
  echo "📦 Installing Air for live reload..."

  if ! curl -sSfL "$AIR_INSTALL_URL" | sh -s; then
    echo "❌ Failed to install Air"
    exit 1
  fi

  if ! ./bin/air -v; then
    echo "❌ Failed to verify Air installation"
    exit 1
  fi

  echo "✅ Air installed successfully"
}

# Setup environment file
setup_env() {
  echo "🔧 Setting up environment file..."

  if [ ! -f "$BASE_ENV_FILE" ]; then
    echo "❌ Error: $BASE_ENV_FILE not found"
    exit 1
  fi

  # Copy base env file
  cp "$BASE_ENV_FILE" "$ENV_FILE"
  echo "✅ Created $ENV_FILE from $BASE_ENV_FILE"

  # Generate secure DB password
  echo "🔐 Generating secure database password..."
  local secure_password=$(openssl rand -base64 32)

  # Update DB_PASSWORD in .env file based on OS type
  if [ "$OS_TYPE" = "macos" ]; then
    sed -i '' "s|^DB_PASSWORD=.*|DB_PASSWORD='$secure_password'|" "$ENV_FILE"
  else
    sed -i "s|^DB_PASSWORD=.*|DB_PASSWORD='$secure_password'|" "$ENV_FILE"
  fi

  echo "✅ Environment file configured successfully"
}

# =============================================
# Main Script
# =============================================

echo "🚀 Starting setup process..."
echo "----------------------------------------"

# Step 1: Get system information
get_system_info

# Display system information
echo -e "\n📋 System Information"
echo "----------------------------------------"
echo "| Operating System    | $OS_TYPE $OS_VERSION"
echo "| Docker Version      | $DOCKER_VERSION"
echo "----------------------------------------"

# Step 2: Install Air
echo -e "\n🛠  Installing Dependencies"
echo "----------------------------------------"
install_air

# Step 3: Setup environment
echo -e "\n🔧  Environment Setup"
echo "----------------------------------------"
setup_env

# Step 4: Completion
echo -e "\n✨ Setup completed successfully!"
echo "🌟 Your development environment is ready to use"
echo "📝 Make sure to review the .env file and adjust any settings if needed"
