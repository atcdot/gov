## GOV: golang version manager
---

 GOV is a simple tool to manage multiple versions of golang. It uses `go` command to install and remove golang versions. It also creates a symlink to the version you want to use.
 Tested on macOS.

## Usage

First of all, you need to install `go` on your system and add the following line to your `~/.bashrc` or `~/.zshrc`:

```bash
source $HOME/.gov
```

### List golang versions (installed and available)

```bash
$ gov list
```

### Install golang version

```bash
$ gov install 1.19.8
```


### Use golang version

```bash
$ gov use 1.19.8
```

### Remove golang version

```bash
$ gov remove 1.19.8
```

### Cleanup golang version alias to use system go

```bash
$ gov cleanup
```

### GUI mode

```bash
$ gov gui
```
