# Hello

A simple web server with some configuration knobs, for
use in practicing configuration.

The server emits a page that contains its version,
followed by a greeting, followed by the value specified
in request path.

A request path of `/quit` exits the server.


### Configuration Knobs

*  Integer `port` flag, default 8080.

*  A hard-coded version const set to zero, change it to make
   ambiguous differences between binaries.

*  Boolean flag `enableRiskyFeature`, default false.
   If enabled, the greeting is italicized.

*  Environment variable `ALT_GREETING`.
   If set, the value overrides the default greeting _Hello_.

Instructions to containerize the hello server are
[here](https://github.com/monopole/hello/blob/master/containerize.md).
