
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
- [Dependencies](https://github.com/egor-denisov/golang-tg-cloud#dependencies)
- [Credits](https://github.com/egor-denisov/golang-tg-cloud#credits)

## Live Demo
Live demo is available here: 
- [Telegram bot](https://t.me/StorageTest1Bot);
- [Web version](https://web-tg-cloud.vercel.app/).
## Final product

At the moment, the project is under development. Now some features of the api of the service and Telegram bot are available for operation. A minimal web interface on react has been developed and deployed.

#### Examples of bot work

<div style="display: flex;">
  <img src="https://sun9-43.userapi.com/impg/BsXAOQ08UIfhqvRvjzrD1nqFPZ-vaA-GuQmXSw/fF1WG-ZMz4Y.jpg?size=591x1280&quality=95&sign=ed309f47873c5c78a81295d6d7993009&type=album" width="300" title="Uploading picture and files"/>
  <img src="https://sun9-27.userapi.com/impg/x51dT-TkouAF5AaobKjCmN08UWR-387fVSLoSA/Xd4W8esq9e4.jpg?size=591x1280&quality=95&sign=b6aa0ec3f96c72b2fb1a02025330ab51&type=album" width="300" title="Getting file"/>
  <img src="https://sun9-32.userapi.com/impg/Fu1q5XLkl5QfFq41B5lfmcOpaDrS2cpxlGHh7g/qX552dJKhgM.jpg?size=591x1280&quality=95&sign=677fa105e374d383e08d3626cfb6ca41&type=album" width="300" title="Creating folder"/>
</div>


#### Examples of web version work

- Login page:
![](https://sun9-33.userapi.com/impg/yxXvwaBgQKcgmni529yaqQEiH9FxQYdHUcE2Lg/7XhlsR535fg.jpg?size=1920x1033&quality=95&sign=fdac6be833147511ee3fda94296dd916&type=album "Login page")

- Folder content:
![](https://sun9-62.userapi.com/impg/iFgJut9DmHAZypvpTMN6hfCYS0kgnNXcUp2LuQ/x7wzbDpgoCI.jpg?size=1920x1033&quality=95&sign=64e29b9a9595d35e27d52b94bcfce334&type=album "Folder content")

- Preview:
![](https://sun9-12.userapi.com/impg/UT3ItuiT8AzmX1eukm91QWYDEy5Z132XLfBIbA/JQl26RdJ_xs.jpg?size=1920x1033&quality=95&sign=d04586815e73872cdeb8df8d9ced4ac1&type=album "Preview")

- Deleting file:
![](https://sun9-22.userapi.com/impg/XSYD0x1HVCaW0m28EhpK8FTEJL7bfzbpLOonKg/gH1rGwm9zKk.jpg?size=1920x1033&quality=95&sign=5961ba1aaad2674e8d1a460965789c08&type=album "Deleting file")

- Uploading files:
![](https://sun9-72.userapi.com/impg/okfS5Br5cVIgm_yrP0FoPbjDRM5xwWdYK_y3jQ/LdMrYY0y18I.jpg?size=1920x1033&quality=95&sign=22ec5b64234498ddf96a5ac7b3dcf8f5&type=album")

## Running the project
1) For starting Telegram version you need get token from https://t.me/BotFather. 
2) To run this project, you will need to add the following environment variables to your ***.env*** file:

`TG_BOT_TOKEN` - Telegram token from BotFather

`DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_NAME` - Database variables.

3) Finally, in cmd run this command:

```
go run main.go
```
Аfter executing, you can write /start to your bot. Also you can build client server with a user-friendly interface.


## API Reference

#### Get info about directory

```http
  GET /directory
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `number` | **Required**. Directory ID |

#### Get info about file

```http
  GET /fileInfo
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

#### Authorization in the system

```http
  GET /auth
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `user_id` | `number` | **Required**. Telegram user_id |
| `username` | `string` | Telegram username |
| `first_name` | `string` | Telegram first name |
| `last_name` | `string` | Telegram lsat name |

#### Create new directory

```http
  GET /createDirectory
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `parent_id` | `number` | **Required**. Id of parent directory |
| `name` | `string` | **Required**. Name of a new directory |
| `user_id` | `number` | **Required**. Telegram user_id |

#### Delete item

```http
  GET /delete
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `number` | **Required**. Item id for deleting |
| `directory_id` | `number` | **Required**. Directory ID |
| `type` | `"file" | "directory"` | **Required**. Type item for deleting |
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

## Dependencies
- Golang v1.20
- PostgreSQL v15.1

## Credits
- [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [postgres](https://github.com/lib/pq)
- [gin](https://github.com/gin-gonic/gin)
- [uuid](https://github.com/google/uuid)
- Backend deployed on [qovery.com](https://qovery.com)
- Database deployed on [clever-cloud.com](https://www.clever-cloud.com/)
