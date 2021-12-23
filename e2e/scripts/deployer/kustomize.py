class KustomizeDeployer():
    def __init__(self, shell):
        self.shell = shell

    def deploy(self, dryrun=False):
        render_command = ["kustomize", "build", "kustomize", "--load_restrictor", "LoadRestrictionsNone", "-o", "rendered.yaml"]
        apply_command = ["kubectl", "apply", "-f", "rendered.yaml"]

        if dryrun:
            apply_command.extend(["--dry-run=server"])
            
        self.shell.execute(render_command)
        self.shell.execute(apply_command)
