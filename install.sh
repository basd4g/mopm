#!/bin/sh
echo "Install mopm..."

if [ "$(whoami)" != "root" ]; then
  echo "Need root privilege"
  exit 1
fi


if [ "$(uname)" = "Linux" ]; then
  target="linux"
elif [ "$(uname)" = "Darwin" ]; then
  target="darwin"
else
  echo "Unsupport Operating System... only linux or macOS"
  exit 1
fi

install_path="/usr/local/bin/mopm"

curl -sL "https://github.com/basd4g/mopm/releases/download/0.0.2/mopm-amd64-${target}.out" > "${install_path}"
chmod a+x "${install_path}"

if which mopm; then
  echo "Mopm is installed!  If you want to uninstall mopm, please delete '${install_path}.'"
else
  "Installing mopm is failed... The script tried to install mopm to '${install_path}', but mopm is not found with 'which mopm'."
fi
