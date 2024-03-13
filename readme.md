# Sun-Short-Link (短链接)


## Usage

### Parameter
```
  -c string
        Path to config file
  -i    Generate example config file
```

### Run
```go
go run main.go -i

go run main.go -c config.yml
```

### Docker Run 

Create file `sun-short-link.yml`.

```yml
port: 8080
links:
  - name: "abc"
    url: "http://example.cc"
    type: 302
  - name: "abcd"
    url: "http://example1.cc"
    type: 302

```

Run
```sh
docker run -p 8080:8080 -v ./sun-short-link.yml:/app/sun-short-link.yml --name sun-short-link hslr/sun-short-link
```