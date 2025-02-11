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
## ECS
![eclogin_v0 0 19](https://github.com/user-attachments/assets/35e3c4fa-cfb3-49db-bb2e-89e85d8d969e)

```
$ eclogin ecs                                                                          
✔ Please enter AWS region: ap-northeast-1
✔ Please enter AWS profile (optional): 
✔ test-cluster
✔ test
✔ xxxxxxxx
✔ test-container
✔ /bin/sh
eclogin equivalent command:
eclogin ecs --cluster test-cluster --task-id xxxxxxxx --container test-container --shell /bin/sh --region ap-northeast-1

If you are using awscli, please copy the following:
aws ecs execute-command \
        --cluster test-cluster \
        --task xxxxxxxx \
        --container test-container \
        --interactive \
        --command /bin/sh \
        --region ap-northeast-1


Starting session with SessionId: ecs-execute-command-xxxxxxxx
# 
```

## EC2
```
$ eclogin ec2
✔ Please enter AWS region (default: ap-northeast-1): █p-northeast-1
Please enter AWS profile (optional): 
✔ test(i-xxxxxxxx)
eclogin equivalent command:
eclogin ec2 --instance-id i-xxxxxxxx --region ap-northeast-1

If you are using awscli, please copy the following:
aws ssm start-session \
        --target i-xxxxxxxx \
        --region ap-northeast-1


Starting session with SessionId: user-xxxxxxxx
sh-4.2$ 
```

## Local
```
$ eclogin local                                                                        
✔ stupefied_dirac(nginx)
✔ /bin/sh
# 
```