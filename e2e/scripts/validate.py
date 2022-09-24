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
              help="Number of validate retry is trying to check",
              default=1,
              show_default=True,
              required=True)
@click.option('--gateway', is_flag=True, help="Validate a gateway")
@click.option("--ratelimited", is_flag=True, help="if it's true, script will check the endpoint should be ratelimited, if it's false, script will check the endpoint should not be ratelimited")
def main(domain, path, retry, gateway, ratelimited):
    shell = ShellDeployer()
    ratelimit_validator = RatelimitValidator(shell, gateway)
    ratelimit_validator.validate(domain, path, retry, ratelimited)


if __name__ == "__main__":
    main()
