# Shell Completion

Shell completion generates a shell completion script for your shell environment. This assumes you have your shell compeltion package installed (`bash-completion` on most RHEL/CentOS systems)

Example:

```shell
source <(helpernodectl completion bash)
```

Usage:

```shell
helpernodectl completion [bash|zsh|fish|powershell]
```

Currently, only the following shells are supported

* bash
* zsh
* fish
* powershell
