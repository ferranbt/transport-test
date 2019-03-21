package main

import (
	"crypto/rand"
	"net"
	"testing"
)

var pipe = pipeTLS

func BenchmarkSendRecv512B(b *testing.B) {
	b.ReportAllocs()
	const payloadSize = 512
	runBenchmark(b, payloadSize)
}

func BenchmarkSendRecv1KB(b *testing.B) {
	b.ReportAllocs()
	const payloadSize = 1024
	runBenchmark(b, payloadSize)
}

func BenchmarkSendRecv256KB(b *testing.B) {
	b.ReportAllocs()
	const payloadSize = 256 * 1024
	runBenchmark(b, payloadSize)
}

func BenchmarkSendRecv512KB(b *testing.B) {
	b.ReportAllocs()
	const payloadSize = 512 * 1024
	runBenchmark(b, payloadSize)
}

func BenchmarkSendRecv1MB(b *testing.B) {
	b.ReportAllocs()
	const payloadSize = 1024 * 1024
	runBenchmark(b, payloadSize)
}

func BenchmarkSendRecv2MB(b *testing.B) {
	b.ReportAllocs()
	const payloadSize = 1024 * 1024 * 2
	runBenchmark(b, payloadSize)
}

func runBenchmark(b *testing.B, size int) {
	b.Run("tls", func(b *testing.B) {
		benchmarkSendRecvConn(b, size, pipeTLS)
	})
	b.Run("rlpx", func(b *testing.B) {
		benchmarkSendRecvRlpx(b, size)
	})
}

func benchmarkSendRecvRlpx(b *testing.B, size int) {
	b.SetBytes(int64(size))
	b.ReportAllocs()
	b.ResetTimer()

	conn0, conn1 := pipeRLPX()

	sendBuf := make([]byte, size)
	rand.Read(sendBuf)

	go func() {
		for i := 0; i < b.N; i++ {
			if _, err := conn0.ReadMsg(); err != nil {
				panic(err)
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		if err := conn1.WriteMsg(0x1, sendBuf); err != nil {
			panic(err)
		}
	}
}

func benchmarkSendRecvConn(b *testing.B, size int, pipe func() (net.Conn, net.Conn)) {
	b.SetBytes(int64(size))
	b.ReportAllocs()
	b.ResetTimer()

	conn0, conn1 := pipe()

	sendBuf := make([]byte, size)
	rand.Read(sendBuf)

	go func() {
		recvBuf := make([]byte, size)
		for i := 0; i < b.N; i++ {
			bytesRead := 0
			for bytesRead < size {
				n, err := conn0.Read(recvBuf)
				if err != nil {
					b.Fatalf("err: %v", err)
				}
				bytesRead += n
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		if _, err := conn1.Write(sendBuf); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}
