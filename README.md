# wait-for-it

[![Build Status](https://travis-ci.org/NoUseFreak/wait-for-it.svg?branch=master)](https://travis-ci.org/NoUseFreak/wait-for-it)

# Install

```bash
curl -L -o /usr/local/bin/wait-for-it https://github.com/NoUseFreak/wait-for-it/releases/download/0.0.6/`uname`_wait-for-it
chmod +x /usr/local/bin/wait-for-it
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
$ wait-for-it
```

# Help

```bash
wait-for-it --help
```
