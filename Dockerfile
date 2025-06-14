# Build-Stage
FROM golang:1.24.3-alpine AS build

# Arbeite im Container im Verzeichnis /app
WORKDIR /app

# Kopiere alle Projektdateien ins Image
COPY . .

# Um statisch zu kompilieren (keine glibc-Abhängigkeit)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Baue die Go-App
RUN go build -o moonapp .

# Finales schlankes Image (scratch = leeres Image ohne OS!)
FROM scratch

# Kopiere Binary und HTML-Template aus dem Build-Image
COPY --from=build /app/moonapp /moonapp
COPY --from=build /app/template.html /template.html

# Port setzen (nur für Doku-Zwecke)
ENV PORT=8082
EXPOSE 8082

# Start-Befehl
ENTRYPOINT ["/moonapp"]
