import sys
import subprocess

class ShellDeployer():
    def __init__(self):
        pass

    def execute(self, command):
        try:
            process = subprocess.Popen(
                command, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
            print(process.stdout.read().strip().decode("utf-8"))
        except subprocess.CalledProcessError as exc:
            print("[ERROR] Command \"{}\" failed with exit code {}: {}".format(
                ' '.join(command), exc.returncode, exc.output.strip().decode("utf-8")))
            sys.exit(1)
        except Exception as exc:
            print(exc.output)
            raise exc
