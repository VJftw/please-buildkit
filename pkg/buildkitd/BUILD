go_library(
    name = "buildkitd",
    srcs = [
        "pleaseworker.go",
        "provider.go",
        "provider-chain.go",
        "provider-docker.go",
        "provider-podman.go",
    ],
    visibility = ["//cmd/..."],
    deps = [
        "//third_party/go/github.com/rs/zerolog",
        "//third_party/go/google.golang.org/grpc",
    ],
)
