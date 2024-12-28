FROM golang:1.23 AS compiling_stage
WORKDIR /go/src/pipeline
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pipeline .

FROM scratch
LABEL ver="1.0"
LABEL maintainer="zombiehunter"
WORKDIR /root/
COPY --from=compiling_stage /go/src/pipeline/pipeline .
CMD ["./pipeline",  "-delay", "20s", "-size", "128", "-log", "console" ]