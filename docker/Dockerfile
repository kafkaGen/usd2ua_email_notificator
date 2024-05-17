FROM golang:1.18

RUN apt-get update && apt-get install -y \
    && apt-get clean

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod tidy

COPY . /app/

CMD ["go", "run", "main.go"]