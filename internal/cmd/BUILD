subinclude("///go//build_defs:go")

go_library(
    name = "cmd",
    srcs = [
        "logging.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go/github.com/mattn/go-isatty",
        "//third_party/go/github.com/rs/zerolog",
    ],
)
