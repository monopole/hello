# Hello

A webserver in Go with some configuration knobs.
For use in practicing configuration.

Knobs:

*  A hard-coded version const set to zero.

*  Boolean flag `enableRiskyFeature`, default false.

*  Integer `port` flag, default 8080.

*  Environment variable `ALT_GREETING`.
   If set, the value overrides the default _Hello_.


Example build, with a version change:

```
package=github.com/monopole/hello
version=2

newSrc=${GOPATH}/src/${package}_${version}.go
newBin=${GOPATH}/bin/${package}_${version}.go

go get -d $package

cat ${GOPATH}/src/${package}.go |\
    sed 's/version = 0/version = '${version}'/' \
    >$newSrc

CGO_ENABLED=0 GOOS=linux \
    go build -o $TUT_IMG_PATH -a -installsuffix cgo $newSrc

go build -o hello_${version}
```

Run and quit:

```
# Start server
ALT_GREETING=salutations \
    $TUT_IMG_PATH --enableRiskyFeature --port 8100 &

# Let it get ready
sleep 2

# Dump html to stdout
curl --fail --silent -m 1 localhost:8100/godzilla

# Send query of death
curl --fail --silent -m 1 localhost:8100/quit
```
