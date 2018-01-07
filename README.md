# Simple Key-Value database with http-endpoints


This is a learning project which implements a simple Key-Value database.


### Learning goals

* Get to know golang
* learn about the stdlib packages
  * http (handlers, middleware)
  * json (un)marshalling
  * testing package
* learn about protocol buffers [google protocol buffers](https://developers.google.com/protocol-buffers/)
* learn about golang concurrency with channels
* learn about database storages


### Implementation

The data storage implementation is based on the descriptions in the book
[Designing Data-Intensive Applications](https://dataintensive.net) chapter 3
with an in-memory hash index.

Currenty supported features:

* SET (HTTP POST)
  * reject requests with wrong Content-Type, application/octet-stream must be set
* GET (HTTP GET)
* DELETE (HTTP DELETE)
* recovery (after restarting the server, the hash-index is rebuild from the underlying file-system)
* kubernetes-ready:
  * dockerized with health and readiness http-endpoints
  * graceful http server shutdown
* prefers json but also supports binary data
* thread-safe (one single write goroutine is used, so database-file is not corrupted during parallel writes)


ToDos:

* use multiple db-log-files
* add "compaction" background job in order to free memory
* use sparse index like SSTables or LSM-Trees
* use snapshots for faster recovery/startup time


### Usage

[httpie](https://httpie.org/) is used as http-client.

```bash
# compile and start static binary (see below how to compile)
PORT=8080 DB_FILENAME=${DB_FILE} ./app

# SET key=mykey value={"foo": "bar"}
http --verbose POST "http://localhost:8080/db/mykey" Content-Type:application/octet-stream foo=bar

# GET
http --verbose GET "http://localhost:8080/db/mykey"

# GET with json-format (http response content-type is set to application/json), use only if you know you stored json!
http --verbose GET "http://localhost:8080/db/mykey?format=json"

# DELETE
http --verbose DELETE "http://localhost:8080/db/mykey"
```


### Development

```bash
# use provided Makefile

# build and run  (build, set env-vars like PORT and DB_FILENAME and runs the server)
make run

# recreate protocol buffer files
make proto

# clean (delete compiled files etc.)
make clean

# run without makefile (assumes static binary is already there)
PORT=8080 DB_FILENAME=${DB_FILE} ./app
```
