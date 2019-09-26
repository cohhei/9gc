9gc: 9gc.go
	go build -o 9gc 9gc.go

test: 9gc
	./test.sh
	go test -v ./*.go

clean:
	rm -f 9gc *.o *~ tmp*

.PHONY: test clean