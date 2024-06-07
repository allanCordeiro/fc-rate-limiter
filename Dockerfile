FROM golang:1.22.3 as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /sample_http cmd/api/main.go


FROM scratch
COPY --from=builder /sample_http /sample_http
COPY --from=builder /usr/src/app/.env .env

EXPOSE 8080
ENTRYPOINT [ "/sample_http" ]