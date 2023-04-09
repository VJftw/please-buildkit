genrule(
    name = "env",
    data = [
        "//third_party/binary/GoogleCloudPlatform/docker-credential-gcr:docker-credential-gcr",
    ],
    outs = [
        "env",
    ],
    cmd = """
set -Eeuo pipefail
cat <<EOF > $OUTS
export PATH="\\\$PATH:\\\$(realpath $(dirname $(out_exe //third_party/binary/GoogleCloudPlatform/docker-credential-gcr:docker-credential-gcr)))"
echo "added docker-credential-gcr to PATH"
EOF
    """,
    visibility = ["PUBLIC"],
)

sh_cmd(
    name = "login",
    data = [
        "//auth/gcp/artifact-registry:env",
        "//third_party/binary/GoogleCloudPlatform/docker-credential-gcr:docker-credential-gcr",
    ],
    shell = "/usr/bin/env bash",
    cmd = """
set -Eeuo pipefail
. $(out_location //auth/gcp/artifact-registry:env)
$(out_exe //third_party/binary/GoogleCloudPlatform/docker-credential-gcr:docker-credential-gcr) \\
    configure-docker -include-artifact-registry
    """,
)
