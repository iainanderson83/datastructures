#!/bin/sh

go test -bench Slice -count 8 > slice_bench.txt

go test -bench StructReturn -count 8 > struct_bench.txt

go test -bench Sprintf -count 8 > sprintf_bench.txt

go test -bench SliceConversion -count 8 > slice_size_bench.txt