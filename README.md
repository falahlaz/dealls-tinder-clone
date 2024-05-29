
# Dealls Technical Interview

In essence, this application is a clone of the well-known dating app Tinder. It employs the Redis geospatial mechanism to offer recommendations for potential dating partners.

## Tech Stack

- Golang
- PostgreSQL
- Redis

## Run Locally

Before you run the project, ensure that you already turn on this services on you local machine:

- PostgreSQL
- Redis

Clone the project

```bash
  git clone https://github.com/falahlaz/dealls-tinder-clone
```

Go to the project directory

```bash
  cd dealls-tinder-clone
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run main.go
```

## Authors

- [@falahlaz](https://www.github.com/falahlaz)

