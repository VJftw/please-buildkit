subinclude("//build/defs:buildkit", "///go//build_defs:go")

buildkit_distroless_image(
    name = "image",
    srcs = [":main"],
    entrypoint = ["/test/scratch1/main"],
)

go_binary(
    name = "main",
    srcs = ["main.go"],
    static = True,
)

buildkit_distroless_image(
    name = "image_with_ca",
    srcs = [":main"],
    entrypoint = ["/test/scratch1/main"],
)

buildkit_distroless_image(
    name = "image_auto_entrypoint",
    srcs = [":main", "data.txt"],
)
