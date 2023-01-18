# ghostmon

**ghostmon** allows us to keep an eye on any currently running gh-ost procceses which are exposed via TCP or Unix Socket.

## Service Configuration
| Config Item            | Env Var           | Type   | Default          | Description                                                                |
|------------------------|-------------------|--------|------------------|----------------------------------------------------------------------------|
| HTTP Address           | `HTTP_ADDR`       | string | `:8080`          | Address we want our HTTP server to run at, e.g `:8080`                     |
| gh-ost Connection Type | `CONNECTION_TYPE` | string | `tcp`            | Connection type we want to use to connect to an external gh-ost process    |
| gh-ost Connection Addr | `CONNECTION_ADDR` | string | `localhost:9001` | Connection address we want to use to connect to an external gh-ost process |
