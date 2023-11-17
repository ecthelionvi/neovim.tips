import subprocess

def run_command(command, cwd=None, background=False):
    if background:
        subprocess.Popen(command, cwd=cwd, shell=True)
    else:
        subprocess.run(command, cwd=cwd, shell=True)


run_command("./neovim-tips", cwd="/neovim-tips/backend", background=True)

run_command("npm start", cwd="/neovim-tips/frontend/neovim-tips", background=True)

run_command("nginx -g 'daemon off;'")
