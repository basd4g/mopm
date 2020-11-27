#!/bin/sh
echo "Install mopm..."

if [ "$(whoami)" != "root" ]; then
  echo "Need root privilege"
  exit 1
fi

if [ "$(uname)" == "Linux" ]; then
  target="linux"
elif [ "$(uname)" == "Darwin" ]; then
  target="darwin"
else
  echo "Unsupport Operating System... only linux or macOS"
  exit 1
fi

curl -sL "https://github.com/basd4g/mopm/releases/download/0.0.2/mopm-amd64-${target}.out" > /usr/bin/mopm
chmod a+x /usr/bin/mopm

echo "mopm is installed!"
