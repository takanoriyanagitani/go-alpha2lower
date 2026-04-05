# Simple Benchmark

## OS: macOS 26

- CPU: M3 Max
- RAM: 48 GB

| cli          | data    | user | sys  | cpu  | total | rate        | ratio |
|:------------:|:-------:|:----:|:----:|:----:|:-----:|:-----------:|:------|
| tr           | 160 MiB | 5.41 | 0.04 | 99%  | 5.466 |    29 MiB/s | (1.0) |
| alpha2lower  | 160 MiB | 0.02 | 0.04 | 70%  | 0.086 | 1,860 MiB/s |  64.1 |
| alpha2lower  | 16 GiB  | 1.45 | 1.17 | 54%  | 4.849 | 3,414 MiB/s | 117.7 |

## OS: Ubuntu 22.04

- CPU: Core i7-13700 
- RAM: 64 GB

| cli          | data    | user | sys  | cpu  | total | rate        | ratio |
|:------------:|:-------:|:----:|:----:|:----:|:-----:|:-----------:|:------|
| tr           | 22 GiB  | 4.65 | 2.51 | 88%  | 8.137 | 2,845 MiB/s | (1.0) |
| alpha2lower  | 22 GiB  | 1.74 | 1.94 | 76%  | 4.844 | 4,778 MiB/s |  1.7  |
