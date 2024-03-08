# About this repo

This repo contains many samples about goroutines, channels and two main concurrency patterns
- Fan-in / Fan-out
- Pipelines

It contains a sample http server where samples can be run.

# Run server

```bash
o run math-server/main.go math-server/controllers.go math-server/server.go    
```

# Testing the math server

```bash
 curl -X POST 'localhost:9999/operate' -H "Content-Type: application/json" --data '{"operation": "multiply", "a": 2, "b": 4}'
# Response {"result":8}
```

# Run samples

After you run the server, you can run each sample like:

```bash
go run 01-ctx-time-out/main.go 
go run 03-pipe-lines/main.go 03-pipe-lines/plot.go # some samples needs more than one file to get it working
```

*Thanks for visiting this repo 🫶*