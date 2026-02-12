# GATOR THE BLOG AGGREGATOR 🐊
**-------------------------------------------------------**

[RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator in Go! We'll call it "Gator", you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## SETUP
**-------------------------------------------------------**

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
**-------------------------------------------------------**

```bash
go install github.com/nico4565/gator@latest
psql -c "CREATE DATABASE gator;"
```

## Config
**-------------------------------------------------------**

To use gator, you must first create a `.gatorconfig.json` within your home directory with the following minimal config (url needs tweaking to work with your configuration of postgres)
```json
{
    "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable"
}

if it doesn't work you can try to remove the ssl part
