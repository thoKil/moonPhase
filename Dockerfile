FROM golang:1.24.3 AS build

WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go build -o moonapp .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=build /app/moonapp .
COPY template.html .

ENV PORT=8082
EXPOSE 8082

CMD ["./moonapp"]