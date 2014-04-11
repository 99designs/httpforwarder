# HTTP forwarder

An agent for forwarding fire-and-forget HTTP requests asynchronously.

The agent acts like a HTTP forward proxy, accepting requests
and proxying them on. However the initial proxy request will always
immediately return with an empty 200 response, and the forwarded request
will be made asyncronously.

The environment variable `LISTEN` listens on a given address. The
address may be in the form of HOST:PORT or PATH, where HOST:PORT is
taken to mean a TCP socket and PATH is a path to a UNIX domain socket.
Defaults to ":8080".

## Example

A web application may need to fire information at an HTTP endpoint, and prefer
a fast, reliable response rather than verification that the HTTP request
succeeded.

For example, while logging an error or some system metrics to a third-party
HTTP service, the application may want to fire-and-forget the HTTP request, and
not be hindered by the chance of connection errors, slow response, or non-200
OK.
