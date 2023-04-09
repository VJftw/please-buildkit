VERSION = "2.1.8"

remote_file(
    name = "docker-credential-gcr",
    binary = True,
    exported_files = ["docker-credential-gcr"],
    extract = True,
    url = f"https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v{VERSION}/docker-credential-gcr_{CONFIG.OS}_{CONFIG.ARCH}-{VERSION}.tar.gz",
    visibility = ["PUBLIC"],
)
