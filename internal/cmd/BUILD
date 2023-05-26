subinclude("///go//build_defs:go")

go_library(
    name = "cmd",
    srcs = [
        "logging.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_rs_zerolog//:zerolog",
        "///third_party/go/github.com_rs_zerolog//log",
    ],
)
