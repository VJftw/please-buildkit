subinclude("///go//build_defs:go")

go_library(
    name = "buildkitd",
    srcs = [
        "provider.go",
        "provider-chain.go",
        "provider-podman.go",
        "provider-root-docker.go",
        "provider-rootless-docker.go",
    ],
    visibility = ["//cmd/..."],
    deps = [
        "///third_party/go/github.com_rs_zerolog//:zerolog",
        "///third_party/go/github.com_rs_zerolog//log",
        "///third_party/go/github.com_gofrs_flock//:flock",
    ],
)
