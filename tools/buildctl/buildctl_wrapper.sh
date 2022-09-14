#!/usr/bin/env bash
set -Exeuo pipefail

reporoot="$(pwd | sed 's#/plz-out/.*##g')"
mkdir -p "${reporoot}/plz-out/buildkit"
temp_dir="$(mktemp -d -p "${reporoot}/plz-out/buildkit" pb.XXXXX)"
export XDG_RUNTIME_DIR="$(mktemp -d -p "/run/user/$(id -u)" pb-xdg.XXXXX)"
function cleanup {
    rm -rf "${temp_dir}" || true
    rm -rf "$XDG_RUNTIME_DIR" || true
}
trap cleanup EXIT

mkdir -p "${temp_dir}/rootlesskit"
mkdir -p "${temp_dir}/buildkitd"

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
ROOTLESSKIT="$SCRIPT_DIR/third_party/binary/moby/buildkit/rootlesskit"
BUILDKIT="$SCRIPT_DIR/third_party/binary/moby/buildkit/bin"
BUILDCTL_DAEMONLESS="$SCRIPT_DIR/third_party/binary/moby/buildkit/buildctl-daemonless.sh"

export BUILDCTL="$BUILDKIT/buildctl"
# isolate network from host w/ rootlesskit
export ROOTLESSKIT="$ROOTLESSKIT --state-dir ${temp_dir}/rootlesskit --net=slirp4netns --copy-up=/etc --disable-host-loopback"
export BUILDKITD="$BUILDKIT/buildkitd --oci-worker=true --containerd-worker=false --root ${temp_dir}/buildkitd --rootless --oci-worker-rootless --oci-worker-gc --oci-worker-gc-keepstorage 0"
"$BUILDCTL_DAEMONLESS" "$@"
