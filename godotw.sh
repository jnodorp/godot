#!/usr/bin/env sh

# The GitHub organization.
ORG=jschlichtholz

# The GitHub repository name.
NAME=godot

# Make sure curl is available.
if [ ! -x "$(command -v curl)" ]; then
	echo "curl is required for this wrapper script to run."
	exit 1
fi

# Create configuration directory if required.
CONFIG_DIR="$HOME/.$NAME"
if [ ! -d "$CONFIG_DIR" ]; then
  mkdir "$CONFIG_DIR"
fi

# Create binary directory if required.
BINARY_DIR="$CONFIG_DIR/bin"
if [ ! -d "$BINARY_DIR" ];then
  mkdir "$BINARY_DIR"
fi

# Determine the latest release.
release=`curl --silent "https://api.github.com/repos/$ORG/$NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'`
echo "Latest relase is $release"

# Create release directory if required.
RELEASE_DIR="$BINARY_DIR/$release"
if [ ! -d "$RELEASE_DIR" ]; then
  mkdir "$RELEASE_DIR"
fi

# Determine the OS (and convert to lowercase).
os=`uname -s | tr '[A-Z]' '[a-z]'`

# Download latest release if required.
BINARY="$RELEASE_DIR/$NAME"
if [ ! -f "$BINARY" ]; then
  echo "Downloading godot $release for $os..."
  curl --silent -L https://github.com/$ORG/$NAME/releases/download/$release/$NAME.$os.amd64 -o $BINARY
fi

$BINARY "$@"
