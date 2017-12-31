# Fossil

Fossil is a proxy for Graphite.

Why would I need a proxy?

Well if you're confortable with the idea of pushing data in clear over the wire without cipher them, ok then. But it's an issue for us :smile:.

Fossil is a drop-down replacement for your Graphite deployment. It will listen TCP connection on the port 2003 and translate them into [sensision](http://www.warp10.io/getting-started/#data-format) before flush them in a directory.

We use it in combination with [Beamium](https://github.com/ovh/beamium). This combination allow us to send our data in a ciphered way.

## Installation

### Dependencies

Before install fossil you need some tools:

* [glide](http://glide.sh/)
* [Go](https://golang.org/)
* make

### Install from sources

First, you have to clone it:

```sh
git clone https://github.com/ovh/fossil.git $GOPATH/src/github.com/ovh/fossil
```

Go into the fresh installation:

```sh
cd $GOPATH/src/github.com/ovh/fossil
```

Install project dependencies using glide:

```sh
glide install
```

Now, the best part compilation:

```sh
make release
```

Finally, install fossil:

```sh
sudo make install
```

## Usage

```sh
$ fossil --help
Fossil fossil Graphite to beamium forwarder

Usage:
  fossil [flags]
  fossil [command]

Available Commands:
  version     Print the version number

Flags:
  -b, --batch int          batch count per file (default 10000)
      --config string      config file to use
  -d, --directory string   directory to write metrics file (default "./sources")
  -l, --listen string      listen address (default ":2003")
  -t, --timeout int        batch timeout for flushing datapoints (default 5)
  -v, --verbose            verbose output

Use "fossil [command] --help" for more information about a command.
```

## Example

```sh
fossil -d /opt/beamium/source
```

In this example fossil, will listen `TCP` connection on port `2003` and will translate all graphite datapoints into sensision and flush them into `/opt/beamium/source` directory.

The flush is realized all `5 seconds` or all `10000 datapoints`.
