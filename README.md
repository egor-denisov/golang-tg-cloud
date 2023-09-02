
# Cloud storage based on Telegram

App build backend server to cloud storage. The user can work with both the telegram bot and the web version

Code is written in Golang, database is in PostgreSQL. 

## Timelapse
This is my first project after studying Golang on the Stepik platform.

The project was launched in July 2023 and is still in development.

## Content
- [Live Demo](https://github.com/egor-denisov/golang-tg-cloud#live-demo)
- [Final product](https://github.com/egor-denisov/golang-tg-cloud#final-product)
- [Running the project](https://github.com/egor-denisov/golang-tg-cloud#running-the-project)
- [API Reference](https://github.com/egor-denisov/golang-tg-cloud#api-reference)
- [About the app](https://github.com/egor-denisov/golang-tg-cloud#about-the-app)
- [Features](https://github.com/egor-denisov/golang-tg-cloud#features)
- [Dependencies](https://github.com/egor-denisov/golang-tg-cloud#dependencies)
- [Credits](https://github.com/egor-denisov/golang-tg-cloud#credits)

## Live Demo
Live demo is available here: 
- [Telegram bot](https://t.me/StorageTest1Bot);
- Web version [currently unavailable].
## Final product

At the moment, the project is under development. Now some features of the api of the service and Telegram bot are available for operation. To visualize the possibilities, a frontend is being developed on react.js

## Running the project
1) For starting Telegram version you need get token from https://t.me/BotFather. 
2) To run this project, you will need to add the following environment variables to your ***environment/env.go*** file:

`BOT_TOKEN` - Telegram token from BotFather

`STORAGE_TOKEN` - You can use the same token for this variable.

`HOST` `PORT` `USER` `PASSWORD` `DBNAME` - Database variables.

3) To create necessary tables, in your PostgreSQL database run commands from ***sql/query.sql***.
4) Finally, in cmd run this command:

```
go build
```
Аfter executing, you can run main.exe file and write /start to your bot. Also you can build client server.


## API Reference

#### Get info about directory

```http
  GET /directory
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `number` | **Required**. Directory ID |

#### Get file

```http
  GET /file
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | **Required**. File ID |

#### Get thumbnail

```http
  GET /thumbnail
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | **Required**. File ID |

#### Get available items

```http
  GET /available
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `user_id` | `number` | **Required**. Telegram user_id |
| `directory_id` | `number` | **Required**. Directory ID |

#### Upload file to cloud

```http
  POST /upload
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `user_id` | `number` | **Required**. Telegram user_id|
| `directory_id` | `number` | **Required**. Directory ID |
| `file` | `multipart/form-data` | **Required**. File for uploading |


## About the app

This application is designed to save small files without restrictions on the amount of memory, because the capacities of telegram servers are used.

The backend of the application is written in golang, used by postgres as a database. A frontend is currently being developed on react.js. In the future, a utility will be written to download larger files.

## Features

This paragraph contains some of the features of this project:

- Сlient synchronization (telegram bot and web version)
- Lack of permanent data storage due to the use of Telegram servers
- Api client that allows you to conveniently access the database

## Dependencies
- Golang v1.20
- PostgreSQL v15.1

## Credits
- [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [postgres](https://github.com/lib/pq)
- [gin](https://github.com/gin-gonic/gin)
- [uuid](https://github.com/google/uuid)
