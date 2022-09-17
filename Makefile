build:
	go build
	chmod 775 ./chat

clean:
	rm -f ./chat

all:
	build
