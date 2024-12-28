# Pipeline

## Usage

| Command           | Description                                |
| ----------------- | ------------------------------------------ |
| -h                | display help message                       |
| -delay (duration) | buffer delay (default 15s)                 |
| -size (int)       | buffer size (default 24)                   |
| -log (string)     | destination for log output (default "none")|

```console
go run . -delay 20s -size 128 -log console
```

## Dockerizing

```console
sudo docker build -t pipeline .
sudo docker run -it --name pipeline pipeline
```
