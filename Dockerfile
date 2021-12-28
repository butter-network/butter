FROM golang

WORKDIR /butter
EXPOSE 0-9999
COPY . .
RUN go build ./examples/chat/main.go

CMD ["/butter/global-chat"]