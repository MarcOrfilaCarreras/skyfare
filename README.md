# Skyfare

Skyfare is a Go-based flight data scraper and aggregator designed to collect airline routes, airports, and flight pricing from multiple sources, with built-in SQLite database migration and storage capabilities.

## Features

- Scrapes airline data, including airports, routes, and flights (currently supports Vueling).
- Stores data in a SQLite database with a migration system.
- Repository pattern for clean data access and updates.
- Centralized logging with configurable verbosity.
- Handles API token refresh and session initialization.
- Designed for easy extension to support other airlines or data sources.

## Installation

``` bash
git clone https://github.com/MarcOrfilaCarreras/skyfare.git
cd skyfare
make all
```

## License

See the [LICENSE.md](LICENSE.md) file for details.