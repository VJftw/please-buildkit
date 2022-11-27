subinclude("///go//build_defs:go")

go_library(
    name = "image",
    srcs = [
        "pusher.go",
        "replace.go",
        "repotag.go",
    ],
    visibility = ["//cmd/..."],
    deps = [
        "//third_party/go/github.com/rs/zerolog",
    ],
)

go_test(
    name = "image_test",
    srcs = [
        "pusher_test.go",
        "replace_test.go",
        "repotag_test.go",
    ],
    external = True,
    deps = [
        ":image",
        "//third_party/go/github.com/stretchr/testify",
    ],
)
