# gocache
### Installation
gocache is installed in two parts, installing the service(s), timer and configuring log rotation, and installing gocache by building the binary and configuring gocache via `gocache.cfg.json`.

### Binary
To provide ease of execution, gocache can can be compiled and ran as a binary. This provides easier deployment and doesn't allow tampering.

To build the binary for gocache, execute `go build -o gocache` in the root `gocache/` directory. Once built, the gocache binary only requires its configuration file `gocache.cfg.json`.

## Configuration
To configure gocache, a `gocache.cfg.json` file needs to be placed in base `gocache/` directory (same directory as binary). 
The template `gocache_template.cfg.json` should be used to create suitable `gocache.cfg.json` for each appliance.
```json
{}
```

### Logging
gocache has builtin logging. The directory of the logs is defined by `gocache.cfg.json` and the `stdout`/`stderr` are both appended to a file called `gocache.log` in this directory. 
at the begginning of each line; standard log messages are denoted by `INFO:`, warnings (not fatal errors) are denoted by `WARNING:` and errors are denoted by `ERROR:`

#### Log rotation
gocache does not rotate its own logs, this is scheduled with `logrotate` using the configuration `config/gocache.conf`. Please view `config/gocache.conf` for more information on the configuration
of this projects log rotation.

### Systemd service(s) and log rotate timer
To ensure maintainability gocache operates as a service. This service can be installed by following these instructions:
Please make sure that you are in the directory which you cloned this repository, or modify the commands to suit the paths.
1. Install gocache
    ```
    This section needs some work, currently deciding on creating the require debian install files
    ```
1. Configure gocache
    ```
    sudo nano /etc/gocache/gocache.cfg.json
    ```

## Testing
Testing can be done using the created unit tests. These unit tests should be places in the `gocache` directory with the rest of the package `.go` files.
To test building the package and the below command can be used:
```
sudo ./debian/rules clean build
```
OR if you just want to run unit tests, the following is adequate
```
cd gocache && go test -v
```

## Building
To build a package you will need a few dependencies, these two scripts allow you to run `mk-build-deps`, which locally install the packages dependencies, and `dpkg-buildpackage` to build the package:
```
    sudo apt-get install dpkg-dev devscripts --yes
```
You can then install the dependencies and then build the debian package.
```
    sudo mk-build-deps --install --tool='apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends --yes' debian/control
    sudo dpkg-buildpackage -B -tc
```
For all go packages we currently only support binary-only building, this is because we don't need go development to occur on the test/production servers.

### Installing
When you get to the point of deploying the package, use `dpkg` to install and manage the packages:
```
    sudo dpkg -i <package>
```