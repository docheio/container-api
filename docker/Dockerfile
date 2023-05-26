FROM golang:latest as golang
WORKDIR /root
COPY ./run.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go build run.go

FROM archlinux:latest
WORKDIR /root
RUN  pacman -Syyu --noconfirm unzip
COPY --from=golang /root/run ./
COPY main.sh ./
COPY backup.sh ./
RUN chmod 777 ./run ./main.sh
CMD [ "/bin/bash","-c","/root/main.sh" ]