VERSION = "0.11.0"

BUILD_ARCHES = {
    "amd64": "x86_64",
}

BUILD_ARCH = BUILD_ARCHES[CONFIG.ARCH]

BUILD_OSES = {
    "linux": "Linux",
}

BUILD_OS = BUILD_OSES[CONFIG.OS]

remote_file(
    name = "crane",
    binary = True,
    exported_files = ["crane"],
    extract = True,
    url = f"https://github.com/google/go-containerregistry/releases/download/v{VERSION}/go-containerregistry_{BUILD_OS}_{BUILD_ARCH}.tar.gz",
    visibility = ["PUBLIC"],
)
