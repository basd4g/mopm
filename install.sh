#!/bin/sh
echo "Install mopm..."

if [ "$(whoami)" != "root" ]; then
  echo "Need root privilege"
  exit 1
fi

curl -sL https://github.com/basd4g/mopm/releases/download/0.0.2/mopm-amd64-linux.out > /usr/bin/mopm
chmod a+x /usr/bin/mopm

echo "mopm is installed!"
