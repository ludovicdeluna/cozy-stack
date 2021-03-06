How to install Cozy-stack?
==========================

## Dependencies

- A reverse-proxy (nginx, caddy, haproxy, etc.)
* A SMTP server
* CouchDB 2.0.0

To install CouchDB 2.0.0 through Docker, take a look at our [Docker specific documentation](docker.md).


## Install the binary

You can either download the binary or compile it.

### Download an official release

You can download a `cozy-stack` binary from our official releases:
https://github.com/cozy/cozy-stack/releases. It is a just a single executable
file (choose the one for your platform). Rename it to cozy-stack, give it the
executable bit (`chmod +x cozy-stack`) and put it in your `$PATH`. `cozy-stack
version` should show you the version if every thing is right.

### Using `go`

[Install go](https://golang.org/doc/install), version >= 1.7. With `go` installed and configured, you can run the following command:

```
go get -u github.com/cozy/cozy-stack
```

This will fetch the sources in `$GOPATH/src/github.com/cozy/cozy-stack` and build a binary in `$GOPATH/bin/cozy-stack`.

Don't forget to add your `$GOPATH` to your `$PATH` in your `*rc` file so that you can execute the binary without entering its full path.

```
export PATH="$(go env GOPATH)/bin:$PATH"
```


## Add an instance for testing

You can configure your `cozy-stack` using a configuration file or different
comand line arguments. Assuming CouchDB is installed and running on default
port `5984`, you can start the server:

```bash
cozy-stack serve
```

And then create an instance for development:

```bash
cozy-stack instances add --dev --apps drive,settings,onboarding --passphrase cozy "cozy.tools:8080"
```

The cozy-stack server listens on http://cozy.tools:8080/ by default. See `cozy-stack --help` for more informations.

The above command will create an instance on http://cozy.tools:8080/ with the passphrase `cozy`.

Make sure the full stack is up with:

```bash
curl -H 'Accept: application/json' 'http://cozy.tools:8080/status/'
```

You can then remove your test instance:

```bash
cozy-stack instances rm cozy.tools:8080
```


## And now?

You should configure your DNS, your reverse-proxy, and read [the
documentation](README.md) before creating cozy instances for production.
