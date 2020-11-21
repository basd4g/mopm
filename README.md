# mopm

Mopm (Manager Of Package Manager) is meta package manager for cross platform environment.
Mopm excute package installation commands based on mopm package definition file.

Mopm also records installation date.
It helps you to write dotfiles.

## File description

| path | description |
| ---- | ---- |
| `~/.mopm/` | the directory for mopm |
| `~/.mopm/repos/` | package repository directory |
| `~/.mopm/repos/github.com/basd4g/mopm-defs/definitions/sample.yaml` | definition (\*) |
| `~/.mopm/repos-url` | package repository's url |
| `~/.mopm/stdout.log` | log file |
| `~/.mopm/stderr.log` | log file |

(\*) sample package definition file of the default package repository

## Mopm package definition file

a package is structed by a yaml file.

Mopm package definition file include all the information.

| key | value description | valid string (regex) |
| --- | --- | --- |
| name | the package's name to specify the package | `^[a-z0-9\-]+$` |
| url | the package's project url | `^http(s)?://.+$` |
| description | the package's description | `^.*$` |
| environments[].architecture | target architecture | `^amd64$` |
| environments[].platform | target platform | `^(linux/ubuntu\|darwin)$` |
| environments[].dependencies[] | dependencies' package name to install the package | `^[a-z0-9\-]+$` |
| environments[].privilege | boolean to need root privilege or not  | `^(true\|false)$` |
| environments[].script | installation script for the environment | `^.*$` |

### Samples

```definitions/sample.mopm.yaml
#!/bin/mopm
name: sample
url: https://github.com/basd4g/mopm
description: This is sample package definition script. It cannot be installed.
environments:
  - architecture: amd64
    platform: darwin
    dependencies:
    verification: "false && false"
    privilege: false
    script: |
      echo "This is sample install script. It is no excution anyware."
  - architecture: amd64
    platform: linux/ubuntu
    dependencies:
    verification: "false && false"
    privilege: true
    script: |
      echo "This is sample install script. It is no excution anyware."
```

## TODO

- dependencies chain
- container for testing package definition
- `checkIfError` -> `Exit1IfError`, `checkIfError(errors.New` -> `Exit1()`
