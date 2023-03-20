### Configuration:
Requires a Postgres Database with following environment settings, can be through docker or any other service, .env file is used by the docker compose for a local container
```
POSTGRES_USER=username 
POSTGRES_PASSWORD=password 
POSTGRES_DB=default
HOST=localhost
```
### Usage: 
1. Clone Project and cd in
```bash
$ git clone git@github.com:antonizyla/IMDb6dos.git 
$ cd IMDb6dos/
```

2. Build the Project
```bash
$ go build .
```

3. Seed the database with current IMDb dataset from the internet
```bash
$ ./IMDb6dos --seed
```

4. Run the Program
```bash
$ ./IMDb6dos
```

Alternatively use the Makefile to `Run`, `Seed` and `Clean`
