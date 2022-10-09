sh_cmd(
    name = "buildctl",
    shell = "/usr/bin/env bash",
    cmd = """
set -Exeuo pipefail
reporoot="\\\$(pwd | sed 's#/plz-out/.*##g')"
mkdir -p "\\\${reporoot}/plz-out/buildkit"
temp_dir="\\\$(mktemp -d -p "\\\${reporoot}/plz-out/buildkit" tmp.XXXXX)"
cleanup() {
    chmod -s "\\\${temp_dir}" || true
    chmod -t "\\\${temp_dir}" || true
    chmod 700 "\\\${temp_dir}" || true
    chmod -R 700 "\\\${temp_dir}" || true
    rm -rf "\\\${temp_dir}" || true
}
trap cleanup EXIT

export XDG_RUNTIME_DIR="\\\${temp_dir}/xdg"
mkdir -p "\\\${temp_dir}/xdg"
mkdir -p "\\\${temp_dir}/rootlesskit"
mkdir -p "\\\${temp_dir}/buildkitd"

ROOTLESSKIT="\\\${reporoot}/$(out_exe //third_party/binary/moby/buildkit:rootlesskit)"
BUILDKIT="\\\${reporoot}/$(out_location //third_party/binary/moby/buildkit:buildkit)"
BUILDCTL_DAEMONLESS="\\\${reporoot}/$(out_exe //third_party/binary/moby/buildkit:buildctl-daemonless-script)"

export BUILDCTL="\\\$BUILDKIT/buildctl"
export ROOTLESSKIT="\\\$ROOTLESSKIT --state-dir \\\${temp_dir}/rootlesskit --net=slirp4netns --copy-up=/etc --disable-host-loopback"
export BUILDKITD="\\\$BUILDKIT/buildkitd --oci-worker=true --containerd-worker=false --root \\\${temp_dir}/buildkitd --rootless --oci-worker-rootless --oci-worker-gc --oci-worker-gc-keepstorage 0"
"\\\$BUILDCTL_DAEMONLESS" "\\\$@"
    """,
    data = [
        "//third_party/binary/moby/buildkit:rootlesskit",
        "//third_party/binary/moby/buildkit:buildkit",
        "//third_party/binary/moby/buildkit:buildctl-daemonless-script",
    ],
    visibility = ["PUBLIC"],
)
