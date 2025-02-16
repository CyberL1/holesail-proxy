#!/bin/sh

set -e

if ! command -v unzip > /dev/null; then
  echo "Error: Unzip is required to install Holesail proxy"
  exit 1
fi

dir="$HOME/.holesail-proxy/bin"
zip="$dir/holesail-proxy.zip"
exe="$dir/holesail-proxy"

if [ "$OS" = "Windows_NT" ]; then
  target="windows"
else
  case $(uname -sm) in
  "Darwin x86_64") target="darwin-amd64" ;;
  "Darwin arm64") target="dawin-arm64" ;;
  "Linux aarch64") target="linux-arm64" ;;
  *) target="linux-amd64"
  esac
fi

download_url="https://github.com/CyberL1/holesail-proxy/releases/latest/download/holesail-proxy-${target}.zip"

if [ ! -d $dir ]; then
  mkdir -p $dir
fi

curl --fail --location --progress-bar --output $zip $download_url
unzip -d $dir -o $zip
chmod +x $exe
rm $zip

echo "Holesail proxy was installed to $runtimer_exe"
if command -v holesail-proxy > /dev/null; then
  echo "Run 'holesail-proxy --help' to get started"
else
  case $SHELL in
  /bin/zsh) shell_profile=".zshrc" ;;
  *) shell_profile=".bashrc" ;;
  esac
  echo "export PATH=\"$dir:\$PATH\"" >> $HOME/$shell_profile
fi
