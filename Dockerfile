FROM golang:1.22.1

WORKDIR /app
COPY . /app

CMD ["go", "run", "."]
