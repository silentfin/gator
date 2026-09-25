# gator
a command-line RSS feed aggregator that fetches and stores posts.

## Features
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Prerequisites
- Go
- PostgreSQL


> [!IMPORTANT]
>
> Create this config file in your home directory, `~/.gatorconfig.json` with your PostgreSQL connection string:
>
> ```json
> {
>  "db_url": "protocol://username:password@host:port/database?sslmode=disable"
> }
> ```

## Installation

### Option 1: Install directly
```bash
go install github.com/silentfin/gator@latest
```

### Option 2: Build from source
1. Clone the repository:
```bash
    git clone https://github.com/silentfin/gator.git
```
2. Go to project directory:
```bash
    cd gator
```

3. Build:
```bash
    go build
```

4. Run:
```bash
    ./gator help
```

## Usage

```bash
gator register <name>          # Create a user
gator login <name>             # Set the current user
gator users                    # List all users
gator reset                    # Delete all database data
gator addfeed <name> <url>     # Add and follow a feed
gator feeds                    # List all feeds
gator follow <url>             # Follow an existing feed
gator unfollow <url>           # Unfollow a feed
gator following                # List followed feeds
gator agg <duration>           # Fetch feeds repeatedly, e.g. agg 10s
gator browse [limit]           # Show recent posts (default: 2)
gator help                     # Show this message
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
