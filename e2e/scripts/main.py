#!/usr/bin/env python

import os

from deployer.kustomize import KustomizeDeployer
from deployer.shell import ShellDeployer
from deployer.manifest import ManifestDeployer

DEFAULT_CLUSTER_MANIFEST_DIR = "manifests/"

def main():
    os.chdir("e2e")

    shell = ShellDeployer()

    kustomize = KustomizeDeployer(shell)
    kustomize.deploy()
    
    manifest = ManifestDeployer(shell, DEFAULT_CLUSTER_MANIFEST_DIR)
    manifest.deploy()

if __name__ == '__main__':
    main()
