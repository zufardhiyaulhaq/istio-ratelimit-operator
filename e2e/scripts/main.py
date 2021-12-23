#!/usr/bin/env python

import os
import click

from deployer.kustomize import KustomizeDeployer
from deployer.shell import ShellDeployer
from deployer.manifest import ManifestDeployer

DEFAULT_CLUSTER_MANIFEST_DIR = "manifests/"


@click.command()
@click.option("--usecases",
              help="Usecases directory name",
              required=True)

def main(usecases):
    os.chdir("e2e/usecases/" + usecases)

    shell = ShellDeployer()

    kustomize = KustomizeDeployer(shell)
    kustomize.deploy()
    
    manifest = ManifestDeployer(shell, DEFAULT_CLUSTER_MANIFEST_DIR)
    manifest.deploy()

if __name__ == '__main__':
    main()
