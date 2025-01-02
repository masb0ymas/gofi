set -e

BASE_ENV_FILE=".env.example"
ENV_FILE=".env"

DOCKER_VERSION=$(docker version --format '{{.Server.Version}}')
OS_TYPE=$(grep -w "ID" /etc/os-release | cut -d "=" -f 2 | tr -d '"')
OS_VERSION=$(grep -w "VERSION_ID" /etc/os-release | cut -d "=" -f 2 | tr -d '"')

if [ "$OS_TYPE" = "ubuntu" ]; then
  OS_TYPE="linux"
fi

if [ "$OS_TYPE" = "alpine" ]; then
  OS_TYPE="linux"
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
  OS_TYPE="macos"
fi

echo -e "---------------------------------------------"
echo "| Operating System  | $OS_TYPE $OS_VERSION"
echo "| Docker            | $DOCKER_VERSION"
echo -e "---------------------------------------------\n"

echo " - Initializing runtime with Air Verse"
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s

./bin/air -v

if [ -f "$BASE_ENV_FILE" ]; then
  cp "$BASE_ENV_FILE" "$ENV_FILE"
  echo " - Copy $BASE_ENV_FILE to $ENV_FILE"

  if [ "$OS_TYPE" = "macos" ]; then
    # Generate a secure DB_PASSWORD
    sed -i '' "s|^DB_PASSWORD=.*|DB_PASSWORD='$(openssl rand -base64 32)'|" "$ENV_FILE"
  elif [ "$OS_TYPE" = "linux" ]; then
    # Generate a secure DB_PASSWORD
    sed -i "s|^DB_PASSWORD=.*|DB_PASSWORD='$(openssl rand -base64 32)'|" "$ENV_FILE"
  fi
fi

echo -e "\nYour setup is ready to use!\n"
