Fossil
=========

Fossil is a proxy for Graphite.

Why would I need a proxy?
Well if you're confortable with the idea of pushing data in clear over the wire without SSL, ok then. But it's an issue for us :)

Fossil is a drop-down replacement for your Graphite deployment. It will listen on TCP:2003 and will flush metrics to a directory.
We use it in combination with <a href="https://github.com/ovh/beamium/" target="_blank">Beamium</a> so that the directory where Fossil flushes its data is a source directory for Beamium.


## Dependencies

- glide https://github.com/Masterminds/glide
- Go
- make


## Build

```sh
$ glide install
$ make release
$ ./build/fossil --help
```

## Config

- Listen (optional): (-l :2003)

    Address on which Fossil must listen 2003 is the default Graphite port


- Directory : (-d /to/my/path)

    Directory used to flush metrics, in most cases, you would use a Beamium source path (source-dir)


## example

```sh
$ ./build/fossil -d /opt/beamium/source
```
