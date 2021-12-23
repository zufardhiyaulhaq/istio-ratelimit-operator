class ManifestDeployer():
    def __init__(self, shell, manifest_dir):
        self.shell = shell
        self.manifest_dir = manifest_dir

    def deploy(self, dryrun=False):
        command = ["kubectl", "apply", "-f", self.manifest_dir, "-R"]

        if dryrun:
            command.extend(["--dry-run=server"])

        self.shell.execute(command)
