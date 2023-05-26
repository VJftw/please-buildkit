subinclude("///go//build_defs:go")

go_binary(
    name = "please_buildkit",
    srcs = [
        "build.go",
        "main.go",
        "push.go",
        "replace.go",
        "buildkitd_worker.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//internal/cmd",
        "//pkg/buildkitd",
        "//pkg/image",
        "///third_party/go/github.com_rs_zerolog//:zerolog",
        "///third_party/go/github.com_rs_zerolog//log",
        "///third_party/go/github.com_urfave_cli_v2//:v2",
        "///third_party/go/github.com_avast_retry-go_v4//:v4",
    ],
    static = True,
)

subinclude("//build/defs:buildkit")

_buildkit_tool()
