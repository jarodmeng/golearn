build=go build

squaring:
	$(build) -o squaring gen.go sq.go 1_squaring.go

fanin:
	$(build) -o fanin gen.go sq.go 2_merge.go 2_fanin.go

cancellation:
	$(build) -o cancellation gen.go sq.go 3_merge_cancellation.go 3_cancellation.go

done:
	$(build) -o done gen.go 4_sq_done.go 4_merge_done.go 4_done.go
