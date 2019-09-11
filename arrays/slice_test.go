package main

import "testing"

func BenchmarkSlicePointers(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		slice := make([]*A, 0, 100)
		for j := 0; j < 100; j++ {
			slice = append(slice, &A{B: j, C: j + 1})
		}
	}
}

func BenchmarkSliceNoPointers(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		slice := make([]A, 0, 100)
		for j := 0; j < 100; j++ {
			slice = append(slice, A{B: j, C: j + 1})
		}
	}
}

func BenchmarkSliceHybrid(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		slice := make([]A, 0, 100)
		for j := 0; j < 100; j++ {
			slice = append(slice, A{B: j, C: j + 1})
		}

		slicep := make([]*A, len(slice))
		for j := range slice {
			slicep[j] = &slice[j]
		}
	}
}
