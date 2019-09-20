9gc: 9gc.go
	go build -o 9gc 9gc.go

test: 9gc
	./test.sh

clean:
	rm -f 9gc *.o *~ tmp*

.PHONY: test clean