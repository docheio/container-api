NAME	= container-api


all:
	go build -o $(NAME)

clean:
	rm -f $(NAME)

re: clean all