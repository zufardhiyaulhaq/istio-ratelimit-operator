#!/usr/bin/env python

import click

from deployer.shell import ShellDeployer
from validator.ratelimit import RatelimitValidator
    
@click.command()
@click.option("--domain",
              help="Domain for ratelimit testing",
              required=True)
@click.option("--path",
              help="Path for ratelimit testing",
              required=True)
@click.option("--retry",
              help="Number of validate retry is trying to check if it's rate limited",
              default=1,
              show_default=True,
              required=True)
@click.option('--gateway', is_flag=True, help="Validate a gateway")
def main(domain, path, retry, gateway):
    shell = ShellDeployer()
    ratelimit_validator = RatelimitValidator(shell, gateway)
    ratelimit_validator.validate(domain, path, retry)

if __name__ == "__main__":
  main()

