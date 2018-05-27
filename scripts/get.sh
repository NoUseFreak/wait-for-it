#!/usr/bin/env bash

# Usage: curl https://raw.githubusercontent.com/NoUseFreak/wait-for-it/master/scripts/get.sh | bash

get_latest_release() {
  curl --silent "https://api.github.com/repos/NoUseFreak/wait-for-it/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

download() {
  curl -Ls -o /usr/local/bin/wait-for-it https://github.com/NoUseFreak/wait-for-it/releases/download/$1/`uname`_wait-for-it
}

echo "Looking up latest release"
RELEASE=$(get_latest_release)

echo "Downloading package"
$(download $RELEASE)

echo "Making executable"
sudo chmod +x /usr/local/bin/wait-for-it
