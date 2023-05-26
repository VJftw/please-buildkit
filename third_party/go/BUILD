subinclude("///go//build_defs:go")
subinclude("///shell//build_defs:shell")

go_toolchain(
    name = "toolchain",
    version = "1.20.4",
)

sh_cmd(
    name = "gen_modules",
    cmd = """
    out="\\\$(./pleasew query reporoot)/third_party/go/modules.build_defs"
    echo "GO_MODULES = {" > "\\\$out.new"

    ./pleasew run ///go//tools/please_go:bootstrap -- get -m go.mod \\
    | sed 's#go_repo(module=##g' \\
    | sed 's#, version=#: #g' \\
    | sed 's#)#,#g' \\
    | sed 's#^#  #g' \\
    >> "\\\$out.new"

    echo "}" >> "\\\$out.new"
    mv "\\\$out.new" "\\\$out"
    """,
)

export_file(
    name = "modules",
    src = "modules.build_defs",
)

subinclude(":modules")

for mod, ver in GO_MODULES.items():
    go_repo(
        module = mod,
        version = ver,
    )
