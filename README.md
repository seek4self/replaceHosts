# replacer

The tool to replace the local hosts. Regularly pull github hosts from [ineo6/hosts](https://github.com/ineo6/hosts) and update to local hosts

## usage

```bash
The tool to replace the local hosts.
Version: v1.0.0
Usage of ./replacer:
  -H string
        the new host ip for the '-d'(input domain) flag. (default "0")
  -d string
        domain in local hosts. (default "github")
  -dd
        disable domain hosts.
  -i string
        replace interval. example: '1h30m', 'h' for hour, and 'm' for minute. (default "2h")
  -one
        replace github hosts once.
  -v    print version.
```

## changelog

- v1.0.0
  - pull remote github hosts and update local hosts
  - disable hosts
  - update singal domain
