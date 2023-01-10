FROM golang:1.19.4-windowsservercore-ltsc2022

WORKDIR /butter
EXPOSE 0-9999
COPY . .
RUN go build ./examples/chat/main.go

CMD ["/butter/global-chat"]