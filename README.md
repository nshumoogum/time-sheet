time-sheet
==================
A tracking service for when one clocks in or out

### Installation

#### Database
* Run `brew install mongo`
* Run `brew services start mongodb`

### Configuration

| Environment variable      | Default         | Description
| ------------------------- | --------------- | ----------------------------------------
| BIND_ADDR                 | :10000          | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT | 5s              | The graceful shutdown timeout in seconds
| MONGODB_BIND_ADDR         | localhost:27017 | The MongoDB bind address
| MONGODB_DATABASE          | timesheets      | The MongoDB dataset database
| MONGODB_COLLECTION        | timesheets      | MongoDB collection

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

See [LICENSE](LICENSE.md) for details.
