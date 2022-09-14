BUILDKIT_VERSION = "0.10.4"

ROOTLESSKIT_VERSION = "1.0.1"

remote_file(
    name = "buildkit",
    binary = True,
    entry_points = {
        "buildctl": "bin/buildctl",
        "buildkitd": "bin/buildkitd",
    },
    extract = True,
    url = f"https://github.com/moby/buildkit/releases/download/v{BUILDKIT_VERSION}/buildkit-v{BUILDKIT_VERSION}.{CONFIG.OS}-{CONFIG.ARCH}.tar.gz",
    visibility = ["PUBLIC"],
)

remote_file(
    name = "rootlesskit",
    binary = True,
    exported_files = ["rootlesskit"],
    extract = True,
    url = f"https://github.com/rootless-containers/rootlesskit/releases/download/v{ROOTLESSKIT_VERSION}/rootlesskit-x86_64.tar.gz",
    visibility = ["PUBLIC"],
)

remote_file(
    name = "buildctl-daemonless-script",
    binary = True,
    url = f"https://raw.githubusercontent.com/moby/buildkit/v{BUILDKIT_VERSION}/examples/buildctl-daemonless/buildctl-daemonless.sh",
    visibility = ["PUBLIC"],
)
