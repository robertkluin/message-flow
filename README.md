Message Flow
============

Message Flow is a simple, pluggable, service discovery and message routing
system.  It consists of two key components: 1) a basic service
discovery/registry system and 2) a reverse proxy/router that routes messages
according to the registry.

Configuration
-------------

Message Flow maintains two mappings, a service and a client mapping.  The 
client mapping is very simplistic.  The mapping service provides the ability to
answer a few simple questions about a client:

  - for a given version of a service, what server should client messages be
    routed to,
  - to deliver a message to a client, which message-router should accept the
    message.

The server mapping handles storing service registration data.  It provides
three mechanisms to determine which specific server a client message should be
routed to.  These are explained below in priority order:

  - server: all messages will be routed to this URI.
  - registrar: if there is not a server specified for the client sending the
    message, message flow will ask the registrar where messages from this
    client should be routed.
  - server-list: a list of servers which will have messages uniform-randomly
    routed to them.

When using either registrar or server-list mechanisms, message-flow will use
"consistent" routing by default.  That is, all messages from a given client to
a given service will be routed to the same server.  The TTL of the "cached"
target can be set or disabled completely.  Each service may be keyed with a
version, different versions are treated as separate services.


Architecture
------------

Being pluggable many architectures are possible.  Message Flow consists of two
key components: routers and a backend routing table backend.  The front-end
routers are responsible for handling the routing of messages.  Front-ends
connect to the routing backend to determine routing information.  The routing
table backend is a simple datastore for a which a routing table adapter exists.
Two adapters are currently bundled: in memory and etcd.  The in memory adapter
is suitable for a single-node message-flow cluster.  The etcd backend is
suitable for a highly-scalable message-flow cluster.


How to Contribute
-----------------
Any contributions are appreciated.  The basic contribution cycle:

  1. Fork message-flow on github, 
  2. Make your contribution (documentation improvement, bug fix,
     optimization, enhancement, etc...),
  3. Ensure existing unit tests run, the contribution has any new unit tests
     neccisary, has been verified to work, and has relevant documentation,
  4. Submit a PR back to the main message-flow repository.

All contributed code must be licensed under the same license as message-flow.

