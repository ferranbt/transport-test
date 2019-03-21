
# Benchmark devp2p transports

Intructions to create TLS certificates:

https://gist.github.com/jim3ma/00523f865b8801390475c4e2049fe8c3

Run the benchmarks:

```
$ go test ./... -run=XXX --bench=Benchmark
```

```
goos: linux
goarch: amd64
pkg: github.com/ferranbt/transport-test
BenchmarkSendRecv512B/tls-8         	  200000	      8381 ns/op	  61.09 MB/s	      81 B/op	       2 allocs/op
BenchmarkSendRecv512B/rlpx-8        	   50000	     28776 ns/op	  17.79 MB/s	   41874 B/op	      56 allocs/op
BenchmarkSendRecv1KB/tls-8          	  200000	      6986 ns/op	 146.57 MB/s	      81 B/op	       2 allocs/op
BenchmarkSendRecv1KB/rlpx-8         	   50000	     40148 ns/op	  25.51 MB/s	   45015 B/op	      56 allocs/op
BenchmarkSendRecv256KB/tls-8        	    5000	    227443 ns/op	1152.57 MB/s	     413 B/op	       2 allocs/op
BenchmarkSendRecv256KB/rlpx-8       	     500	   3342478 ns/op	  78.43 MB/s	 1378883 B/op	      67 allocs/op
BenchmarkSendRecv512KB/tls-8        	    3000	    396608 ns/op	1321.93 MB/s	     917 B/op	       4 allocs/op
BenchmarkSendRecv512KB/rlpx-8       	     300	   5509282 ns/op	  95.16 MB/s	 2439465 B/op	      76 allocs/op
BenchmarkSendRecv1MB/tls-8          	    2000	    732851 ns/op	1430.82 MB/s	    1944 B/op	       9 allocs/op
BenchmarkSendRecv1MB/rlpx-8         	     100	  10831108 ns/op	  96.81 MB/s	 6141776 B/op	      96 allocs/op
BenchmarkSendRecv2MB/tls-8          	    1000	   1411436 ns/op	1485.83 MB/s	    6057 B/op	      18 allocs/op
BenchmarkSendRecv2MB/rlpx-8         	      50	  21634261 ns/op	  96.94 MB/s	12383980 B/op	     136 allocs/op
PASS
ok  	github.com/ferranbt/transport-test	19.508s
```
