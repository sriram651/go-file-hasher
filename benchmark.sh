#!/bin/bash

# Usage: ./benchmark.sh <directory> [output_file]
# Example: ./benchmark.sh ./local-process-this results.txt

set -e

DIR="${1:?Usage: ./benchmark.sh <directory> [output_file]}"
OUTPUT="${2:-benchmark_$(date +%Y%m%d_%H%M%S).txt}"
RUNS=2
WORKERS=(1 2 4 8 16 32)
BINARY="./hasher"

if [ ! -d "$DIR" ]; then
  echo "Error: directory '$DIR' not found"
  exit 1
fi

if [ ! -f "$BINARY" ]; then
  echo "Building hasher..."
  go build -o "$BINARY" .
fi

FILE_COUNT=$(find "$DIR" -maxdepth 1 -type f | wc -l | xargs)
DIR_SIZE=$(du -sh "$DIR" | cut -f1)

{
  echo "========================================"
  echo " File Hasher Benchmark"
  echo "========================================"
  echo ""
  echo "Date:       $(date)"
  echo "CPU:        $(sysctl -n machdep.cpu.brand_string 2>/dev/null || lscpu 2>/dev/null | grep 'Model name' | cut -d: -f2 | xargs || echo 'unknown')"
  echo "Cores:      $(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 'unknown')"
  echo "Go version: $(go version)"
  echo "OS:         $(uname -mrs)"
  echo ""
  echo "Directory:  $DIR"
  echo "Files:      $FILE_COUNT"
  echo "Total size: $DIR_SIZE"
  echo ""
  echo "========================================"
  echo " Results"
  echo "========================================"
  echo ""

  for w in "${WORKERS[@]}"; do
    echo "--- Workers: $w ---"
    for r in $(seq 1 $RUNS); do
      echo -n "  Run $r: "
      # Use bash's built-in time via a subshell to capture stderr
      { time "$BINARY" -dir "$DIR" -workers "$w" -q ; } 2>&1 | grep real | awk '{print $2}'
    done
    echo ""
  done

  echo "========================================"
  echo " Raw output (with full time stats)"
  echo "========================================"
  echo ""

  for w in "${WORKERS[@]}"; do
    echo "# $w worker(s)"
    for r in $(seq 1 $RUNS); do
      echo "$ time ./hasher -dir $DIR -workers $w -q"
      { time "$BINARY" -dir "$DIR" -workers "$w" -q ; } 2>&1
      echo ""
    done
  done
} 2>&1 | tee "$OUTPUT"

echo ""
echo "Results saved to: $OUTPUT"
