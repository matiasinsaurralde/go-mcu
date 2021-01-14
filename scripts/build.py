#!/usr/bin/python3
import re, subprocess, os
version_file = open('version.go').read()
version = re.findall("version = \"(.*)\"", version_file)[0]

build_path = 'build'

targets = [
	{'GOOS': 'darwin', 'GOARCH': 'amd64'},
	{'GOOS': 'linux', 'GOARCH': '386'},
	{'GOOS': 'linux', 'GOARCH': 'amd64'},
	{'GOOS': 'linux', 'GOARCH': 'arm'},
	{'GOOS': 'windows', 'GOARCH': '386'},
	{'GOOS': 'windows', 'GOARCH': 'amd64'},
	{'GOOS': 'windows', 'GOARCH': 'arm'}
]

env = os.environ.copy()

try:
	files = os.listdir(build_path)
	for f in files:
		full_path = os.path.join(build_path, f)
		os.remove(full_path)
except:
	print("Creating build directory")
	os.mkdir(build_path)

for t in targets:
	out_file = "go-mcu-{0}-{1}-{2}".format(version, t['GOOS'], t['GOARCH'])
	if t['GOOS'] == 'windows':
		out_file = out_file + ".exe"
	print("Building {0}...".format(out_file))
	out_path = os.path.join(build_path, out_file)
	env['GOOS'] = t['GOOS']
	env['GOARCH'] = t['GOARCH']
	subprocess.run(["go", "build", "-o", out_path], env=env)
