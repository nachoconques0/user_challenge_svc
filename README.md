# Faceit Challenge
### Made with <3 by Juan Calcagno AKA Nacho. 
#### I love CS:GO BTW hardcode fan since 1.5 lol


### What I've done
I've built a microservice that is allowing to do a full CRUD of an USER.

### Features
- [x] Create an User
- [x] Update an User nickname
- [x] Delete an user (soft)
- [x] Find users and by country code also
- [x] Emits an event whenevere an actions happens to the User Entity
- [x] Contains a subcriber that will log whenever an event was sent
- [x] It has validations 
- [x] HTTP Endpoints, including a health check
- [x] gRPC Endpoints


### Postman Collection available :white_check_mark:
It is available in the repo.

### How to run it 🙀
1. Have docker in your machine
2. `git clone` this repo
3. Once you are inside the repo
4. Run `docker compose up -d` this will initiate a container with a running postgres DB
5. Run `make mod` so you download needed pkgs
6. Run `make migration-run dir=up` this will run all needed migrations
7. Run `make run` and if all good. Project should be running ready to get some HTTP calls and also gRPC.


### You don't want to run it? 😈
1. Have docker in your machine
2. `git clone` this repo
3. Once you are inside the repo
4 Run `make test` and this will trigger a docker compose file that will spin up a test DB + mgirations and then run all the needed tests. By the time of writing this test are passing lol. 🤞🏼


### Extra thoughts
I know you guys didn't fully asked to create a pubsub with go. At the beginning what I did was just to log that an event was sent in an `emit` function and this fn was being called from the aggregate. But then I've decided to learn and practice a bit. And did the pubsub using GO. I hope it doesn't backfire. 


## HTTP Endpoints
#### Create User `POST /users`
- All fields must be in payload
##### Body
```
{
 "first_name": "nachotest",
 "last_name":"calcagno",
 "nickname": "nacho",
 "password":"123123123",
 "email":"nachotest@gmail.com",
 "country":"UK"
}
```

##### Response 201
```
{
    "id": "7a634e9a-cafa-4fd2-b914-fde26465b3f7",
    "first_name": "nachotest",
    "last_name": "calcagno",
    "nickname": "nacho",
    "email": "nachotest@gmail.com",
    "country": "UK"
}

```
#### Update User `PATCH /users/{id}`
- Only nickname can be updated
##### Body
```
{
 "nickname": "nachofromCSGO"
}
```

##### Response 200
```
{
    "id": "7a634e9a-cafa-4fd2-b914-fde26465b3f7",
    "first_name": "nachotest",
    "last_name": "calcagno",
    "nickname": "nachofromtheUKUpdated",
    "email": "nachotest@gmail.com",
    "country": "UK"
}
```

#### Delete User `DELELTE /users/{id}`
- It does a soft delete
- No body needed
##### Response 200

#### Find Users `GET /users?limit=3&page=1&country=UK`
- it will fail if there are wrong query params
##### Response 200
```
[
    {
        "id": "b9fc0e89-845d-4c2a-a7dd-07a7282da2d6",
        "first_name": "nachoeventtest",
        "last_name": "calcagno",
        "nickname": "nacho",
        "email": "nachoeventtest@gmail.com",
        "country": "UK"
    },
    {
        "id": "316078df-97dc-4615-9601-f004f42c80ec",
        "first_name": "nachoeventtest1",
        "last_name": "calcagno",
        "nickname": "nacho",
        "email": "nachoeventtest1@gmail.com",
        "country": "UK"
    },
    {
        "id": "ffd96b86-0b8c-47b0-82de-9efc4ec4d8d5",
        "first_name": "nachoeventtest11111111111",
        "last_name": "calcagno",
        "nickname": "nacho",
        "email": "nachoeventtest111111111@gmail.com",
        "country": "UK"
    }
]
```





### Project folder structure 🌴
```
📦user_challenge_svc
 ┣ 📂cmd
 ┃ ┗ 📂server
 ┃ ┃ ┣ 📜dev.go
 ┃ ┃ ┗ 📜main.go
 ┣ 📂migrations
 ┃ ┣ 📜20250419101833_init.down.sql
 ┃ ┣ 📜20250419101833_init.up.sql
 ┃ ┣ 📜20250419103243_user-table.down.sql
 ┃ ┣ 📜20250419103243_user-table.up.sql
 ┃ ┣ 📜20250420081608_add-user-event-table.down.sql
 ┃ ┗ 📜20250420081608_add-user-event-table.up.sql
 ┣ 📂pkg
 ┃ ┗ 📂challenge
 ┃ ┃ ┣ 📂app
 ┃ ┃ ┃ ┣ 📜instance.go
 ┃ ┃ ┃ ┗ 📜options.go
 ┃ ┃ ┣ 📂db
 ┃ ┃ ┃ ┣ 📜db.go
 ┃ ┃ ┃ ┗ 📜options.go
 ┃ ┃ ┣ 📂env
 ┃ ┃ ┃ ┗ 📜env.go
 ┃ ┃ ┣ 📂helpers
 ┃ ┃ ┃ ┗ 📜db.go
 ┃ ┃ ┣ 📂internal
 ┃ ┃ ┃ ┣ 📂aggregate
 ┃ ┃ ┃ ┃ ┗ 📂user
 ┃ ┃ ┃ ┃ ┃ ┣ 📜user.go
 ┃ ┃ ┃ ┃ ┃ ┗ 📜user_test.go
 ┃ ┃ ┃ ┣ 📂controller
 ┃ ┃ ┃ ┃ ┣ 📂grpc
 ┃ ┃ ┃ ┃ ┃ ┗ 📂user
 ┃ ┃ ┃ ┃ ┃ ┃ ┣ 📜controller.go
 ┃ ┃ ┃ ┃ ┃ ┃ ┗ 📜controller_test.go
 ┃ ┃ ┃ ┃ ┣ 📂http
 ┃ ┃ ┃ ┃ ┃ ┗ 📂user
 ┃ ┃ ┃ ┃ ┃ ┃ ┣ 📜controller.go
 ┃ ┃ ┃ ┃ ┃ ┃ ┗ 📜controller_test.go
 ┃ ┃ ┃ ┃ ┗ 📂pubsub
 ┃ ┃ ┃ ┃ ┃ ┗ 📜user.go
 ┃ ┃ ┃ ┣ 📂entity
 ┃ ┃ ┃ ┃ ┗ 📂user
 ┃ ┃ ┃ ┃ ┃ ┣ 📂event
 ┃ ┃ ┃ ┃ ┃ ┃ ┗ 📜user.go
 ┃ ┃ ┃ ┃ ┃ ┣ 📜user.go
 ┃ ┃ ┃ ┃ ┃ ┗ 📜user_test.go
 ┃ ┃ ┃ ┣ 📂mocks
 ┃ ┃ ┃ ┃ ┣ 📜mock_publisher.go
 ┃ ┃ ┃ ┃ ┣ 📜mock_user_aggregate.go
 ┃ ┃ ┃ ┃ ┗ 📜mock_user_service.go
 ┃ ┃ ┃ ┣ 📂model
 ┃ ┃ ┃ ┃ ┣ 📜error.go
 ┃ ┃ ┃ ┃ ┗ 📜user.go
 ┃ ┃ ┃ ┣ 📂repo
 ┃ ┃ ┃ ┃ ┣ 📜user.go
 ┃ ┃ ┃ ┃ ┗ 📜user_test.go
 ┃ ┃ ┃ ┗ 📂service
 ┃ ┃ ┃ ┃ ┗ 📂user
 ┃ ┃ ┃ ┃ ┃ ┣ 📜service.go
 ┃ ┃ ┃ ┃ ┃ ┗ 📜service_test.go
 ┃ ┃ ┣ 📂proto
 ┃ ┃ ┃ ┗ 📂user
 ┃ ┃ ┃ ┃ ┣ 📜user.pb.go
 ┃ ┃ ┃ ┃ ┣ 📜user.proto
 ┃ ┃ ┃ ┃ ┗ 📜user_grpc.pb.go
 ┃ ┃ ┣ 📂pubsub
 ┃ ┃ ┃ ┣ 📂local
 ┃ ┃ ┃ ┃ ┗ 📜bus.go
 ┃ ┃ ┃ ┣ 📜publisher.go
 ┃ ┃ ┃ ┗ 📜subscriber.go
 ┃ ┃ ┗ 📂server
 ┃ ┃ ┃ ┣ 📂grpc
 ┃ ┃ ┃ ┃ ┗ 📜grpc.go
 ┃ ┃ ┃ ┗ 📂http
 ┃ ┃ ┃ ┃ ┣ 📂middleware
 ┃ ┃ ┃ ┃ ┃ ┗ 📜traceid.go
 ┃ ┃ ┃ ┃ ┣ 📜http.go
 ┃ ┃ ┃ ┃ ┣ 📜options.go
 ┃ ┃ ┃ ┃ ┗ 📜routes.go
 ┣ 📜.gitignore
 ┣ 📜.golangci.yml
 ┣ 📜Makefile
 ┣ 📜README.md
 ┣ 📜docker-compose.yml
 ┣ 📜docker-compose_test.yml
 ┣ 📜generate-mocks.sh
 ┣ 📜go.mod
 ┗ 📜go.sum
  ```