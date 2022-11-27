subinclude("///go//build_defs:go")

go_binary(
    name = "please_buildkit",
    srcs = [
        "main.go",
        "push.go",
        "replace.go",
        "worker.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//internal/cmd",
        "//pkg/buildkitd",
        "//pkg/image",
        "//third_party/go/github.com/mattn/go-isatty",
        "//third_party/go/github.com/rs/zerolog",
        "//third_party/go/github.com/urfave/cli/v2",
    ],
)

subinclude("//build/defs:buildkit")

_buildkit_tool()
