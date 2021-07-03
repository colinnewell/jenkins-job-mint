# jenkins-job-mint

## Overview

Mint new jenkins jobs quickly.

## Building

	git clone https://github.com/colinnewell/jenkins-job-mint.git
	cd jenkins-job-mint
	make
	sudo make install

## Config file

The tool allows you to provide connection settings in multiple ways, command
line flags, config file, and via env vars.

Config files are probably the most sensible way to store the connection info
for Jenkins:

This is expected to live in: `~/.jenkins-job-mint.yaml`


```yaml
user: admin
token: xxxxxxxxxxxxx
url: http://localhost:8080
```

You can override which config file to load using `--config file`.

Env vars are expected to start `MINT_`.  So `MINT_URL=https://jenkins` will
override the Jenkins URL for example.

## Using

On it's own it provides help and works in the usual way a multi command tool
works.

    mint

To get the Jenkins config for jobs to make it easy for you to create templates.

    mint get -v test-job test-job2

Once you have the config files edit them to replace parts that need changing
when a new job is created.  The `job` variable is provided as standard, others
can be passed in using the `--variables` command line arg.

The template can use the standard Go text/template syntax.

```xml
  ...
  <builders>
    <hudson.tasks.Shell>
      <command>echo test
echo {{ .job | html }} - {{ .a | html }}
</command>
    </hudson.tasks.Shell>
  </builders>
```

Create a new job using a template:

	mint job --template template.xml --variables '{"a":3}' new-job
