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
CONFIGURATION_DIRECTORY="$HOME/.$NAME"
if [ ! -d "$CONFIGURATION_DIRECTORY" ]; then
  mkdir "$CONFIGURATION_DIRECTORY"
fi

# Create binary directory if required.
BINARY_DIRECTORY="$CONFIGURATION_DIRECTORY/bin"
if [ ! -d "$BINARY_DIRECTORY" ];then
  mkdir "$BINARY_DIRECTORY"
fi

# Determine the latest release.
release=`curl --silent "https://api.github.com/repos/$ORG/$NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'`
echo "Latest relase is $release"

# Create release directory if required.
RELEASE_DIRECTORY="$BINARY_DIRECTORY/$release"
if [ ! -d "$RELEASE_DIRECTORY" ]; then
  mkdir "$RELEASE_DIRECTORY"
fi

# Determine the OS (and convert to lowercase).
os=`uname -s | tr '[A-Z]' '[a-z]'`

# Download latest release if required.
BINARY="$RELEASE_DIRECTORY/$NAME"
if [ ! -f "$BINARY" ]; then
  echo "Downloading godot $release for $os..."
  curl --silent -L https://github.com/$ORG/$NAME/releases/download/$release/$NAME.$os.amd64 -o $BINARY
fi

$BINARY "$@"
