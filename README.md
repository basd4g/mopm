# mopm

Mopm (Manager Of Package Manager) is meta package manager for cross platform environment.
Mopm excute package installation commands based on mopm package definition file.

Mopm package definition file includes the following infomation.

- package name
- installation script for each platform
- dependencies to run installation script
- to need root privilege or not to need root privilege

Mopm also records installation date.
It helps you to write dotfiles.

## Mopm package definition file

A package is structed by some files.
A mopm package definition file is define a package on a architecture and a platform.

Mopm package definition file is a bash shell script with specified comments.

### File content

Mopm package definition file follow these rules

#### Line 1 is shebang.

`#!/bin/bash -e`

#### Include all the information

Include following these all information for mopm package.

There are the format of a line including key and value: '#mopm {key}: {value}'

| --- | --- | --- |
| key | value description | valid string (regex) |
| name | the package's name to specify the package | ^[a-z0-9\\-]+$ |
| url | the package's project url | ^http(s).://[ ]+$ |
| description | the package's description | ^.\*$ |
| architecture | target architecture | ^amd64$ |
| platform | target platform | ^(ubuntu|darwin)$ |
| dependencies | dependencies' package name to install the package | ^([a-z0-9\\-]+(, [a-z0-9\\-]+)\*)?$ |
| privilege | to need root privilege or not or never | ^(root|unnecessary|never)$ |

#### Install script

Lines starting without # is installation bash script.

### Samples

```amd64-ubuntu-yarn.mopm
#!/bin/bash -e
# This is mopm package definition file. Please excute on mopm.
#mopm name: yarn
#mopm url: https://classic.yarnpkg.com
#mopm description: Fast, reliable, and secure dependency management.
#mopm architecture: amd64
#mopm platform: ubuntu
#mopm dependencies: curl, apt, apt-key, grep
#mopm verification: "which yarn"
#mopm privilege: root
# privilege ... root, unnecessary, never
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
YARN_PACKAGE_URL='deb https://dl.yarnpkg.com/debian/ stable main'
YARN_LIST='/etc/apt/sources.list.d/yarn.list'
if ! grep -q "$YARN_PACKAGE_URL" "$YARN_LIST" ; then
  echo "$YARN_PACKAGE_URL" >> "$YARN_LIST"
fi
apt update
apt install -y yarn
```

```amd64-darwin-yarn.mopm
#!/bin/bash -e
# This is mopm package definition file. Please excute on mopm.
#mopm name: yarn
#mopm url: https://classic.yarnpkg.com
#mopm description: Fast, reliable, and secure dependency management.
#mopm architecture: amd64
#mopm platform: darwin
#mopm dependencies: brew
#mopm verification: "which yarn"
#mopm privilege: never
brew install yarn
```

