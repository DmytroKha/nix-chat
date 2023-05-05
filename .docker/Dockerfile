FROM golang:1.18-alpine

WORKDIR /app/src

COPY backend/go.mod .
COPY backend/go.sum .

RUN go mod download

COPY backend .

COPY backend/internal/infra/database/migrations ../migrations

RUN go build -o ../nixchat cmd/server/main.go

WORKDIR /app

RUN rm -rf src

EXPOSE 8080

CMD [ "/app/nixchat" ]
#ENTRYPOINT ["tail","-f", "/dev/null"]