FROM golang

WORKDIR /butter
EXPOSE 0-9999
COPY . .
RUN go build ./examples/global-chat/main.go

CMD ["/butter/global-chat"]