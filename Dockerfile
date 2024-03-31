FROM golang:1.22.1

COPY . .

CMD ["go", "run", "."]
