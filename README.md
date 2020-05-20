# TunBot

I made this for a now abandoned project, the idea was that the client would be able to generate TCP tunnels on the fly for RDP connections by making HTTP calls to the server endpoints via a simple API.

The idea was that customers trying to access their RDP sessions in certain countries could reduce their latency significantly if they session was tunneled through certain endpoints.

The idea worked in practice and reduced latencies by 50ms or more.

## notable functionality

* uses the magic packet wizardry to actually detect an RDP server on the target
* does some other stuff including basic HTTP API
* works last time i tried it?
* doesn't have any secrets in the repo that are in use anywhere
