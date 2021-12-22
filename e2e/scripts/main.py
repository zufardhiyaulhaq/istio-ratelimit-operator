#!/usr/bin/env python

import os

from deployer.kustomize import KustomizeDeployer
from deployer.shell import ShellDeployer

def main():
    os.chdir("e2e")

    shell = ShellDeployer()
    kustomize = KustomizeDeployer(shell)
    kustomize.deploy()

if __name__ == '__main__':
    main()
