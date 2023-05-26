BUILDKIT_VERSION = "0.11.6"

remote_file(
    name = "_buildkit#download",
    extract = True,
    url = f"https://github.com/moby/buildkit/releases/download/v{BUILDKIT_VERSION}/buildkit-v{BUILDKIT_VERSION}.{CONFIG.OS}-{CONFIG.ARCH}.tar.gz",
)

genrule(
    name = "buildctl",
    binary = True,
    srcs = [":_buildkit#download"],
    outs = ["buildctl"],
    cmd = "mv $SRCS/buildctl $OUTS",
    visibility = ["PUBLIC"],
)

genrule(
    name = "buildkitd",
    binary = True,
    srcs = [":_buildkit#download"],
    outs = ["buildkitd"],
    cmd = "mv $SRCS/buildkitd $OUTS",
    visibility = ["PUBLIC"],
)
