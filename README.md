# replacer

🌍 *[English](README.md) ∙ [简体中文](README_CN.md)*

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/seek4self/replaceHosts/Go)](https://github.com/seek4self/replaceHosts/actions/workflows/go.yml)
[![GitHub](https://img.shields.io/github/license/seek4self/replaceHosts)](https://github.com/seek4self/replaceHosts/blob/master/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/seek4self/replaceHosts?display_name=tag)](https://github.com/seek4self/replaceHosts/releases)

The tool to replace the local hosts. Regularly pull github hosts from [ineo6/hosts](https://github.com/ineo6/hosts) and update to local hosts

## usage

> **Notice**: Running this software under Windows environment will be blocked by Microsoft Defender or anti-virus software. Please add this software to the trust list or white list.  
> **MUST** be run as administrator or `sudo`

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
  - pull remote github hosts and update local hosts
  - disable hosts
  - update singal domain
