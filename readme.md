# Sun-Short-Link (短链接服务 | 链接中转服务)


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

### 高级功能

#### 1. 支持目标链接关键字替换 示例

配置文件：
```yml
  - name: update_log
    url: https://xxx.top/zh_cn/update/{version}
    replace:
      - query_name: version_name # get参数名字
        keyword: '{version}' # url地址替换的关键字
    
```

请求地址：`https://link.cc/update_log?version_name=1.0.0`

目标地址：`https://xxx.top/zh_cn/update/1.0.0`

#### 2. 支持定义链接组

此功能一般用于多语言链接组功能

配置文件：
```yml
  - name: donate_lang
    type: 302
    urls_query_name: lang
    urls:
      - name: zh_cn
        url: https://github.com/hslr-s/sun-panel?v=v{version}
      - name: en
        url: https://doc.sun-panel.top/introduce/donate.html
    replace:
      - query_name: version_name # get参数名字
        keyword: '{version}' # url地址替换的关键字
```

请求地址：`https://link.cc/donate_lang?version_name=1.7&lang=zh_cn`

目标地址：`https://github.com/hslr-s/sun-panel?v=v{version}`