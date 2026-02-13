# Benchmark: Concurrent File Hashing (SHA-256)

## Environment

| Parameter       | Value                          |
|-----------------|--------------------------------|
| CPU             | Apple M4 (10 cores)            |
| Memory          | 16 GB                          |
| OS              | macOS (Darwin 25.2.0, arm64)   |
| Go version      | go1.25.6 darwin/arm64          |

## Dataset

| Property    | Value                                           |
|-------------|-------------------------------------------------|
| Directory   | `./local-process-this`                          |
| Total files | 185                                             |
| Total size  | ~117 MB                                         |
| File types  | `.txt`, `.bin`, `.dat`, `.log`, `.csv`, `.json`, `.xml`, `.md` |
| File sizes  | Range from a few KB to several MB               |

## Method

Each worker count was run twice using `time` to capture wall-clock duration.
Quiet mode (`-q`) was used to exclude I/O overhead from printing results.

```
time ./hasher -dir ./local-process-this -workers <N> -q
```

## Results

| Workers | Run 1 (s) | Run 2 (s) | Avg (s) | CPU Utilization | Speedup vs 1 Worker |
|--------:|----------:|----------:|--------:|----------------:|--------------------:|
|       1 |     0.063 |     0.062 |  0.0625 |          ~100%  |              1.00x  |
|       2 |     0.042 |     0.037 |  0.0395 |          ~180%  |              1.58x  |
|       4 |     0.023 |     0.022 |  0.0225 |          ~338%  |              2.78x  |
|       8 |     0.014 |     0.017 |  0.0155 |          ~548%  |              4.03x  |
|      16 |     0.012 |     0.015 |  0.0135 |          ~615%  |              4.63x  |
|      32 |     0.014 |     0.016 |  0.0150 |          ~604%  |              4.17x  |

## Raw Output

```
# 1 worker
./hasher -dir ./local-process-this -workers 1 -q  0.05s user 0.01s system 99% cpu 0.063 total
./hasher -dir ./local-process-this -workers 1 -q  0.05s user 0.02s system 101% cpu 0.062 total

# 2 workers
./hasher -dir ./local-process-this -workers 2 -q  0.05s user 0.02s system 173% cpu 0.042 total
./hasher -dir ./local-process-this -workers 2 -q  0.05s user 0.02s system 188% cpu 0.037 total

# 4 workers
./hasher -dir ./local-process-this -workers 4 -q  0.06s user 0.02s system 332% cpu 0.023 total
./hasher -dir ./local-process-this -workers 4 -q  0.06s user 0.02s system 344% cpu 0.022 total

# 8 workers
./hasher -dir ./local-process-this -workers 8 -q  0.06s user 0.02s system 571% cpu 0.014 total
./hasher -dir ./local-process-this -workers 8 -q  0.06s user 0.03s system 526% cpu 0.017 total

# 16 workers
./hasher -dir ./local-process-this -workers 16 -q  0.06s user 0.02s system 647% cpu 0.012 total
./hasher -dir ./local-process-this -workers 16 -q  0.06s user 0.03s system 583% cpu 0.015 total

# 32 workers
./hasher -dir ./local-process-this -workers 32 -q  0.06s user 0.02s system 631% cpu 0.014 total
./hasher -dir ./local-process-this -workers 32 -q  0.06s user 0.03s system 577% cpu 0.016 total
```

## Analysis

- **Linear scaling up to 4 workers**: Going from 1 to 4 workers yields a near-linear 2.78x speedup, indicating the workload is well-distributed and I/O-bound contention is low.
- **Diminishing returns beyond 8 workers**: The jump from 8 to 16 workers provides only a marginal improvement (~0.002s), suggesting the bottleneck shifts from CPU to disk I/O or goroutine scheduling overhead.
- **No benefit at 32 workers**: Performance is effectively the same as 8-16 workers. On a 10-core CPU, exceeding the physical core count adds coordination cost without throughput gains.
- **Sweet spot**: For this dataset and hardware, **8 workers** offers the best balance of speed and resource efficiency.
