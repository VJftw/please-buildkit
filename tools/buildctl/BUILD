sh_binary(
    name = "buildctl",
    main = "buildctl_wrapper.sh",
    deps = [
        "//third_party/binary/moby/buildkit:buildctl-daemonless-script",
        "//third_party/binary/moby/buildkit:buildkit",
        "//third_party/binary/moby/buildkit:rootlesskit",
    ],
    visibility = ["PUBLIC"],
)
