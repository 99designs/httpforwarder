# HTTP forwarder

An agent for forwarding http requests asynchronously.

The agent acts like a HTTP forward proxy, accepting requests
and proxying them on. However the initial proxy request will always
immediately return with an empty 200 response, and the forwarded request
will be made asyncronously.

The environment variable `LISTEN` listens on a given address. The
address may be in the form of HOST:PORT or PATH, where HOST:PORT is
taken to mean a TCP socket and PATH is a path to a UNIX domain socket.
Defaults to ":8080".
