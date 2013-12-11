# Http forwarder

An agent for forwarding http requests asyncronously.

The agent acts like a HTTP forward proxy, accepting requests
and proxying them on. However the initial proxy request will always
immediately return with an empty 200 response, and the forwarded request
will be made asyncronously.
