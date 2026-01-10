module github.com/ngtcp2/treemap-go/treemap/bench

go 1.25

require (
	github.com/emirpasic/gods/v2 v2.0.0-alpha
	github.com/google/btree v1.1.3
	github.com/ngtcp2/treemap-go v0.0.0-00010101000000-000000000000
	github.com/tidwall/btree v1.8.1
)

replace github.com/ngtcp2/treemap-go => ../..
