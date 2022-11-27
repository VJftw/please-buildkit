subinclude("///go//build_defs:go")

go_library(
    name = "buildkitd",
    srcs = [
        "pleaseworker.go",
        "provider.go",
        "provider-chain.go",
        "provider-podman.go",
        "provider-root-docker.go",
        "provider-rootless-docker.go",
    ],
    visibility = ["//cmd/..."],
    deps = [
        "//third_party/go/github.com/rs/zerolog",
        "//third_party/go/google.golang.org/grpc",
    ],
)
