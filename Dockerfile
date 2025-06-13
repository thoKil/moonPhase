FROM golang:1.22 AS build

WORKDIR /app
COPY . .
RUN go build -o moonapp .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=build /app/moonapp .
COPY template.html .

ENV PORT=8082
EXPOSE 8082

CMD ["./moonapp"]
