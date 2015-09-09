#! /bin/bash

set -eu
cd `dirname $0`

NAME="ledcubesim"
VERSION="$(git describe --always --dirty) ($(date --date="@$(git show -s --format='%ct' HEAD)" '+%Y-%m-%d'))"
WORKSPACE="$PWD"
BIN="$WORKSPACE/bin"
LIB="$WORKSPACE/lib"
GOPATH="$LIB/gopath"

mkdir -p "$BIN"
mkdir -p "$LIB"

usage() {
cat << EOF
SYNOPSIS

    This is a very handy tool to manage this Go project.

USAGE
      $ $0 install    # Install dependencies
      $ $0 build      # Build a binary
      $ $0 run        # Run the currently built binary
EOF
}

case ${1:-} in
	"install")
		GPM="$LIB/gpm"
		if [ ! -e "$GPM" ]; then
			echo "*** Installing Dependency Manager ***"
			wget -qO- https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm > "$GPM"
			chmod 755 "$GPM"
		fi

		echo "*** Installing Dependencies ***"
		cd "$WORKSPACE/src"
		GOPATH="$GOPATH" "$GPM" install
		cd "$WORKSPACE"
	;;

	"build")
		echo "*** Building Project ***"
		cd "$WORKSPACE/src"
		GOPATH="$GOPATH" go build -o "$WORKSPACE/bin/$NAME"
		cd "$WORKSPACE"
	;;

	"run")
		"$BIN/$NAME"
	;;

	*)
		usage
	;;
esac
