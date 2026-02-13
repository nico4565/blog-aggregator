# GATOR THE BLOG AGGREGATOR 🐊


[RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator in Go! We'll call it "Gator", you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## SETUP


You will need:

- Postgres v15+:
    - macOS with brew
        ```bash
        brew install postgresql@15

    - Linux / WSL (Debian):
        ```bash
        sudo apt update
        sudo apt install postgresql postgresql-contrib
        ```

        Ensure the installation worked. The psql command-line utility is the default client for Postgres. Use it to make sure you're on version 15+ of Postgres:
        ```bash
        psql --version

    - (Linux only) Update postgres password:
        ```bash
        sudo passwd postgres

- Go toolchain (version 1.25+):
    - The easiest way (in my opinion) is to use [Webi](https://webinstall.dev/golang/) to install Go.
    - Otherwise, the [official download page](https://golang.org/doc/install) is also a great option.


## Installation


```bash
go install github.com/nico4565/gator@latest
psql -c "CREATE DATABASE gator;"
```

## Config


To use gator, you must first create a `.gatorconfig.json` within your home directory with the following minimal config (url needs tweaking to work with your configuration of postgres)
```json
{
    "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable"
}

if it doesn't work you can try to remove the ssl part
```

# Gator CLI Commands

A quick reference for the common commands in the **Gator** CLI tool.  

## Common Commands

| Command | Description |
|---------|-------------|
| `register <name>` | Create a new user |
| `login <name>` | Set the current user |
| `addfeed <name> <url>` | Add a feed with a name and URL |
| `follow <url>` | Follow a feed by URL |
| `unfollow <url>` | Unfollow a feed by URL |
| `agg <duration>` | Continuously fetch posts at the given interval (e.g., `30s`, `5m`) |
| `browse [limit]` | List recent posts, optionally limited |

## Examples

```bash
./gator register alice
./gator login alice
./gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
./gator agg 1m
./gator browse 20

# Development Notes

- The `agg` command is typically run in its own terminal since it loops on an interval.  
- If you see database errors, verify that `DATABASE_URL` is set correctly and that Postgres is reachable.  

## Project Ideas

- Sorting/filtering for `browse`  
- Pagination support  
- Concurrency improvements in `agg`  
- Fuzzy search  
- Bookmarks / likes  
- TUI interface  
- HTTP API with authentication
