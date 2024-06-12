FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -R
RUN CGO_ENABLED=0 GOOS=linux go build -o pwrotate ./cmd/pwrotate

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/pwrotate /
CMD ["/pwrotate"]

