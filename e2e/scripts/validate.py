#!/usr/bin/env python

import click

from validator.ratelimit import RatelimitValidator
    
@click.command()
@click.option("--domain",
              help="Domain for ratelimit testing",
              required=True)
@click.option("--path",
              help="Path for ratelimit testing",
              required=True)
def main(domain, path):
    ratelimit_validator = RatelimitValidator()
    ratelimit_validator.validate(domain, path)

if __name__ == "__main__":
  main()

