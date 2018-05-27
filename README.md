# wait-for-it

[![Build Status](https://travis-ci.org/NoUseFreak/wait-for-it.svg?branch=master)](https://travis-ci.org/NoUseFreak/wait-for-it)

`wait-for-it` is a cli tool to check if a service is available. It was build to run between starting docker services
and running integration tests.

The tool is plugin based, only the plugins for your specific cases are installed.

# Install

```bash
curl -sL http://bit.ly/get-wait-for-it | bash
```

# Usage

Create your config file `wait-for-it.yml`.

```yaml
services:
  mysql_check:
    plugin: mysql
    host: localhost
    parameters:
      port: 3306
      username: root
      password: root
```

Now check if the mysql service is available.

```bash
wait-for-it
```

# Help

```bash
wait-for-it --help
```
