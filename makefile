CC=g++
CFLAGS=-Ofast
LDFLAGS=-lntl

all: ssearch

ssearch: ssearch.cpp
	$(CC) $(CFLAGS) ssearch.cpp -o ssearch $(LDFLAGS)

clean:
	rm ssearch
