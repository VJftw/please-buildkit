VERSION = "0.57.0"

remote_file(
    name = "syft",
    binary = True,
    exported_files = ["syft"],
    extract = True,
    url = f"https://github.com/anchore/syft/releases/download/v{VERSION}/syft_{VERSION}_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    visibility = ["PUBLIC"],
)
