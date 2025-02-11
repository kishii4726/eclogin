# eclogin
CLI tool for logging into AWS EC2/ECS/Local docker containers.

# Requirement
- [session-manager-plugin](https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

# Install
```
$ curl -OL <Release assets url>

$ tar -zxvf <Download file name>

$ sudo mv eclogin /usr/local/bin
```

# Usage
```
$ eclogin ec2
$ eclogin ecs
$ eclogin local
```
