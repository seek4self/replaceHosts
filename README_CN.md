# replacer

🌍 *[English](README.md) ∙ [简体中文](README_CN.md)*

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/seek4self/replaceHosts/Go)](https://github.com/seek4self/replaceHosts/actions/workflows/go.yml)
[![GitHub](https://img.shields.io/github/license/seek4self/replaceHosts)](https://github.com/seek4self/replaceHosts/blob/master/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/seek4self/replaceHosts?display_name=tag)](https://github.com/seek4self/replaceHosts/releases)

一个替换本地 hosts 的工具。 定期从 [ineo6/hosts](https://github.com/ineo6/hosts) 仓库拉取 github hosts 并更新到本地 hosts。

## usage

> **注意**: 在 Windows 环境中运行该软件会被 Microsoft Defender 或者杀毒软件拦截并隔离，请将该软件添加到信任列表或者白名单中。  
> **必须**使用管理员或 `sudo` 身份运行该软件

```text
The tool to replace the local hosts.
Version: v1.0.0
Usage of ./replacer:
  -D string
        domain in local hosts. (default "github")
  -H string
        the new host ip for the '-D'(input domain) flag. (default "0")
  -d    run app as a daemon.
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
  - 拉取远端 github hosts 并更新到本地 hosts
  - 禁用 hosts
  - 更新单个域名 host
