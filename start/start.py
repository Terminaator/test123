import subprocess
import os


def run(command):
    return subprocess.run(command.split(" "), stdout=subprocess.PIPE, cwd=os.getcwd() + "\\start\\docker").stdout.decode('utf-8')


command_containers_down = "docker-compose down"
command_containers_up = "docker-compose up -d"
command_containers_list = "docker container ls"

commands = [command_containers_down,
            command_containers_up,
            command_containers_list]

for command in commands:
    response = run(command)
    print(response)
