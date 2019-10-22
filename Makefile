SRCS=$(wildcard *.go)
OBJS=$(SRCS:.c=.o)

9gc: $(OBJS)
	go build -o 9gc $(SRCS)

test: 9gc
	./test.sh
	go test -v ./*.go

clean:
	rm -f 9gc *.o *~ tmp*

.PHONY: test clean