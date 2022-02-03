# Go Deer - Core


* [Getting Started](#Getting-Started)
* [Components](#Components)
* [Help and Support](#Help-and-Support)
* [Migration Guide](#Migration-Guide)
* [Contributing](#Contributing)

## Getting Started

By default, `go get` will bring in the latest tagged release version of the module.

```shell
go get github.com/mercadolibre/fury_go-core
```

To get a specific release version of the module use `@<tag>` in your `go get` command.

```shell
go get github.com/mercadolibre/fury_go-core@v1.0.0
```

To get the latest module repository change use `@latest`.

```shell
go get github.com/mercadolibre/fury_go-core@latest
```

To run all the tests of this SDK use the included `Makefile`.

```shell
make
```

## Components

### [log](./pkg/log)

Package `log` uses ZAP fmt, which is a small wrapper around [Uber log package](https://godoc.org/go.uber.org/zap).




> Please don't hesitate to ask! Your question will likely be useful for other people and help us to improve the existing documentation.


