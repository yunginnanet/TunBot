# TunBot

I made this many moons ago for a now abandoned project, the idea was that the client would be able to generate TCP tunnels on the fly for RDP connections by making HTTP calls to the server endpoints via a simple API.

The idea was that VPS customers trying to access their servers via RDP in certain countries could reduce their latency significantly if the connection to the RDP server was tunneled through specific endpoints.

The idea worked in practice and reduced latencies by 50ms or more.

## Notes

* Uses the magic packet wizardry to actually detect an RDP server on the target
* Basic HTTP API on server will allow for multiple tunnels to exist within one program instance
* Works last time i tried it? (tested in win10 and debian)
* Doesn't have any secrets in the repo that are in use anywhere

## To-do

#### Client

* Calculate the best server endpoint by running tcping latency checks
* Web interface

#### Server

* SPLIT UP SERVER INTO ONE MASTER SERVER WITH ADDITIONAL SATELLITE SERVERS
  * Master server holds DB of satellite server info
  * Master server collects logs from nodes (central logging)
  * Health checks on satellite servers
  * Satellite server installation script

* Administration Panel

* Statistics
   * Currently active tunnels/nodes
   * Bandwidth statistics per node

* Automatic cleanup for defunct/old tunnels

#### Overall

* General security auditing, this thing is definitely not ready for prod
