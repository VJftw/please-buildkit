VERSION = "1.6"

BUILD_ARCHES = {
    "amd64": "64",
}

BUILD_ARCH = BUILD_ARCHES[CONFIG.ARCH]

remote_file(
    name = "jq",
    url = f"https://github.com/stedolan/jq/releases/download/jq-{VERSION}/jq-{CONFIG.OS}{BUILD_ARCH}",
    visibility = ["PUBLIC"],
    binary = True,
)
