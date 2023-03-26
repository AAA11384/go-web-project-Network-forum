# go-web-project-Network-forum

a simple net forum developed by go language

The project finished by 2023.3.26 aim to practice basic go language and usage of modules, and cost me about two weeks. Guiding by Wenzhou Li (Teacher Qimi).

## Intro

### start

1.Set proper configuration in `settings/config.yaml`, and execute `create_table.sql`in your databases.

2.start the program using `config.yaml` by `go run main.go settings/config.yaml`.you can change parameter of path of configuration profile when your configuration profile in other path.

Warning : This project has no front-end files, please use the interface test program to use.

### Functions

It is a simple forum provide functions of post community, creating a post, creating user, voting for a post and return posts in specific order.

###  API

#### sign up

`POST` `127.0.0.1:8081/api/v1/signup`

Example Body

```json
{
    "username": "aaa11384",
    "password": "12345678",
    "re_password": "12345678",    
    "email": "12356@qwe.com"
}
```

#### log in

`POST` `127.0.0.1:8081/api/v1/login`

Example Body

```js
{
    "username": "aaa1138499002",
    "password": "12345678"
}
```

#### post

`POST` `127.0.0.1:8081/api/v1/post` Remember to carry Token !!!

Example Body

```json
{
    "title": "Test post",
    "content": "i love go",
    "community_id": 1
}
```

#### vote

`POST` `127.0.0.1:8081/api/v1/vote/`

Example Body

```json
{
    "post_id": "162998689709690880",
    "direction": "1"
}
```

#### list all community

`GET` `127.0.0.1:8081/api/v1/community`

#### Inquiry information of a community

`GET` `127.0.0.1:8081/api/v1/community/communityid`

Example `127.0.0.1:8081/api/v1/community/1`

#### Inquiry information of a post

`GET` `127.0.0.1:8081/api/v1/post/POSTID`

Example `127.0.0.1:8081/api/v1/post/159386697560231936`

#### Inquiry information of all post by page

`POST` `127.0.0.1:8081/api/v1/posts?page=1&size=10`

#### Inquiry information of all post by page and order

`GET` `127.0.0.1:8081/api/v1/community_posts?page=1&pageSize=10&CommunityID=0&order=time`

when `CommunityID=0` search for all posts in order

#### Inquiry information of post in specific community by order

`GET` `127.0.0.1:8081/api/v1/community_posts?page=1&pageSize=10&CommunityID=1&order=time`

## Mods

|    Mods    |                           Function                           |                         Address                          |
| :--------: | :----------------------------------------------------------: | :------------------------------------------------------: |
|    Gin     |                route and register middleware                 |      [click here](https://github.com/gin-gonic/gin)      |
|   Viper    |   read all the config in one profile, convenient to change   |       [click here](https://github.com/spf13/viper)       |
|    Zap     |            record the log to debug and supervise             |       [click here](https://github.com/uber-go/zap)       |
|    Sqlx    | offering interfaces connecting to Mysql database by sql language |      [click here](https://github.com/jmoiron/sqlx)       |
|  go-redis  |       offering interfaces connecting to Mysql database       |     [click here](https://github.com/redis/go-redis)      |
| validator  | verify number and format of parameters form user and fill into a model struct | [click here](https://github.com/go-playground/validator) |
| Ratelimit  |      control the ability of accepting request of server      |    [click here](https://github.com/uber-go/ratelimit)    |
| snowflake  |   generate unique and increasing ID information for users    |   [click here](https://github.com/bwmarrin/snowflake)    |
| golang-JWT | generate token for user and verify token from user to relieve the pressure of server |     [click here](https://github.com/golang-jwt/jwt)      |

## Highlight spot

### Three-level construction 

Controllers response for verify parameters, return error messages and call the methods in logic.

Logic will split a request into many parts and execute them one by one, methods in logic will call methods in dao if the request related to databases.

Dao contains methods related to change and inquire mysql and redis, to reduce program coupling, the methods always small.

### Return values and error management

The http code of information returned by server for all requests is 200, if there is a error happened, the return value contains statute code and error information, The 404 page is called by the front-end processing return information.

There are consistent error format in controllers and logic/mysql, every kind of errors has it own code. Error values become clear and low  coupling.

### Data cache

There are two keys in redis : to store ID of all the posts, `bluebell:post:time` is in order of time,`bluebell:post:score`is in order of score. Every communities has many posts, so every community has a key in redis, the key is `bluebell:community:communityid`, value is ID of all the posts in this community.

There is an API, let user inquire all posts in one community by specified order, we use `ZInterStore` to generate a new key named `community/communityid/order` and set 60s expire time to it. Expire time let server don't need to execute `ZInterStore` every time when a request coming. Because it is enormous cost when there are huge posts in server.

### Redis Txpipeline

When a post is created, we need to change keys of community and post in two different orders. It must success or fail at the same time. When key `community/communityid/order` was created, we also need to set Expire time of it at the same time. The program use Txpipeline not only solve the problem but also accelerate speed.

### graceful shut down

The program use os.flag to receive system signal, when a signal coming, invoke Shutdown function and set the max waiting time of five seconds. The Shutdown function can close free connection and waiting for active connection, to develop user experience.











