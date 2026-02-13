# go-file-hasher

A concurrent file hasher written in Go. Walks a directory, computes SHA-256 hashes for every file using a configurable worker pool.

## Build

```bash
go build -o hasher .
```

## Usage

```bash
./hasher -dir <path> [-workers <n>] [-q]
```

| Flag       | Default | Description                          |
|------------|---------|--------------------------------------|
| `-dir`     | —       | **(required)** Directory to process  |
| `-workers` | 5       | Number of concurrent goroutines      |
| `-q`       | false   | Quiet mode — only print path + hash  |

### Examples

```bash
# Hash all files with default 5 workers
./hasher -dir ./my-files

# Use 8 workers, quiet output
./hasher -dir ./my-files -workers 8 -q
```

## Benchmarking

A benchmark script is included to test performance across different worker counts.

```bash
# Run against any directory
./benchmark.sh ./my-files

# Custom output file
./benchmark.sh ./my-files results.txt
```

See [BENCHMARK.md](BENCHMARK.md) for sample results on Apple M4.

## How it works

1. `filepath.WalkDir` traverses the target directory (skips hidden dirs)
2. File paths are sent to a buffered channel
3. Worker goroutines pull paths off the channel and compute SHA-256 hashes
4. A `sync.WaitGroup` ensures all files are processed before exit
