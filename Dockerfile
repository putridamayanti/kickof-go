FROM golang:1.22-alpine

WORKDIR /app
COPY . .
RUN go mod tidy

COPY *.go ./

RUN go build -o /kickof-go

EXPOSE 80

CMD [ "/kickof-go" ]