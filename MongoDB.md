
# Start mongodb by docker

## pull image 

```shell
$ docker pull mongo:latest
```

## start container

```shell
$ docker run -itd -p 27017:27017 -v /data/mongodb:/data/mongodb --name mongodb mongo:latest
```

## create auth user and password 

```shell
$ docker exec -it mongodb mongo admin

# create user named 'admin' and password is '123456'
>  db.createUser({ user:'admin',pwd:'123456',roles:[ { role:'userAdminAnyDatabase', db: 'admin'},"readWriteAnyDatabase"]});
# try auth
> db.auth('admin', '123456')
```

