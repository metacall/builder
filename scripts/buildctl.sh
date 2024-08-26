#!/bin/sh
# buildctl-daemonless.sh spawns ephemeral buildkitd for executing buildctl.
#
# Usage: buildctl-daemonless.sh build ...
#
# Flags for buildkitd can be specified as $BUILDKITD_FLAGS .
#
# The script is compatible with BusyBox shell.
set -eu

: ${BUILDCTL=buildctl}
: ${BUILDCTL_CONNECT_RETRIES_MAX=10}
: ${BUILDKITD=buildkitd}
: ${BUILDKITD_FLAGS=}
: ${ROOTLESSKIT=rootlesskit}
: ${LIMACTL=limactl}

# $tmp holds the following files:
# * pid
# * addr
# * log
tmp=$(mktemp -d /tmp/buildctl-daemonless.XXXXXX)

if [ "$(uname)" = "Darwin" ]; then
    trap "rm -rf $tmp" EXIT
else
    trap "kill \$(cat $tmp/pid) || true; wait \$(cat $tmp/pid) || true; rm -rf $tmp" EXIT
fi

startBuildkitd() {
    addr=
    helper=
    if [ "$(uname)" = "Darwin" ]; then
        echo "MacOS: $(sw_vers -productName) $(sw_vers -productVersion) detected"
        if which buildctl >/dev/null 2>&1; then
            echo "buildctl is installed at $(which buildctl)"
            if which lima > /dev/null 2>&1; then
                echo "lima is installed at $(which lima)"
                limactlStart
                addr="unix://$HOME/.lima/buildkit/sock/buildkitd.sock"
                echo "$addr" > "$tmp/addr"
                return
            else
                echo "Please install lima for running buildkitd using : brew install lima"
                exit 1
            fi
        else
            echo "builtctl is not installed. Please install it via brew install buildkit or build it using the moby/buildkit repo."
            exit 1
        fi
    elif [ $(id -u) = 0 ]; then
        addr=unix:///run/buildkit/buildkitd.sock
    else
        addr=unix://$XDG_RUNTIME_DIR/buildkit/buildkitd.sock
        helper=$ROOTLESSKIT
    fi
    $helper $BUILDKITD $BUILDKITD_FLAGS --addr=$addr >$tmp/log 2>&1 &
    pid=$!
    echo $pid >$tmp/pid
    echo $addr >$tmp/addr
}

limactlStart() {
    if ! $LIMACTL list | grep -q buildkit; then
        echo "Instance not up, running instance..."
        $LIMACTL start template://buildkit --tty=false
    else
        echo "Lima: Buildkit instance is already up"
    fi
}

# buildkitd supports NOTIFY_SOCKET but as far as we know, there is no easy way
# to wait for NOTIFY_SOCKET activation using busybox-builtin commands...
waitForBuildkitd() {
    addr=$(cat $tmp/addr)
    try=0
    max=$BUILDCTL_CONNECT_RETRIES_MAX
    until $BUILDCTL --addr=$addr debug workers >/dev/null 2>&1; do
        if [ $try -gt $max ]; then
            echo >&2 "could not connect to $addr after $max trials"
            echo >&2 "========== log =========="
            cat >&2 $tmp/log
            exit 1
        fi
        sleep $(awk "BEGIN{print (100 + $try * 20) * 0.001}")
        try=$(expr $try + 1)
    done
}

startBuildkitd
waitForBuildkitd
$BUILDCTL --addr=$(cat $tmp/addr) "$@"
