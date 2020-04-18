<h1 align="center">watchexec - Watch for Changes and Execute Commands</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/watchexec"><img src="https://goreportcard.com/badge/github.com/cbergoon/watchexec?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/watchexec"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

```watchexec``` watches a directory for changes, executing the commands in the configuration file. 

#### Install
```
$ go get github.com/cbergoon/watchexec
$ go install github.com/cbergoon/watchexec
```

#### Example Usage
`configuration.yaml`
```yaml
ignore: [hello, main]
commands:
  - executable: go
    arguments: [version]
    sequence: 10
  - executable: go
    arguments: [build, hello.go]
    sequence: 10
```

```sh
    $ watchexec
```
#### License
This project is licensed under the MIT License.