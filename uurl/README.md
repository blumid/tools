# uurl

### use this tool in order to make your urls unique according to query string.

* example:

```bash
▶ cat urls.txt
https://abc.com?a=1&b=2
https://abc.com?a=2=333&b=9
https://abc.com/home?a=10&b=5

▶ cat urls.txt | uurl
https://abc.com?a=1&b=2
https://abc.com/home?a=10&b=5
```

* Usage example:

```bash
cat urls | uurl
```

* Install:

```
go install github.com/blumid/tools/uurl@latest
```
