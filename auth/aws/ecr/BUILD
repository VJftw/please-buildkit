genrule(
    name = "env",
    data = [
        "//third_party/binary/awslabs/amazon-ecr-credential-helper:docker-credential-ecr-login",
    ],
    outs = [
        "env",
    ],
    cmd = """
cat <<EOF > $OUTS
export PATH="\\\$PATH:\\\$(realpath $(dirname $(out_exe //third_party/binary/awslabs/amazon-ecr-credential-helper:docker-credential-ecr-login)))"
echo "added docker-credential-ecr-login to PATH"
EOF
    """,
    visibility = ["PUBLIC"],
)

sh_cmd(
    name = "login",
    data = [
        "//auth/aws/ecr:env",
        "//third_party/binary/awslabs/amazon-ecr-credential-helper:docker-credential-ecr-login",
    ],
    cmd = """
set -e
. $(out_location //auth/aws/ecr:env)
$(out_exe //third_party/binary/awslabs/amazon-ecr-credential-helper:docker-credential-ecr-login) \\
    --help
    """,
)
