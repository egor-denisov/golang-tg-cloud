
# Cloud storage based on Telegram

App build backend server to cloud storage. The user can work with both the telegram bot and the web version

Code is written in Golang, database is in PostgreSQL. 

## Timelapse
This is my first project after studying Golang on the Steppik platform.

The project was launched in July 2023 and is still in development.

## Content
- [Live Demo](https://github.com/egor-denisov/golang-tg-cloud#live-demo)
- [Final product](https://github.com/egor-denisov/golang-tg-cloud#final-product)
- [Running the project](https://github.com/egor-denisov/golang-tg-cloud#running-the-project)
- [About the game](https://github.com/egor-denisov/golang-tg-cloud#about-the-game)
- [Features](https://github.com/egor-denisov/golang-tg-cloud#features)
- [Dependencies](https://github.com/egor-denisov/golang-tg-cloud#dependencies)
- [Credits](https://github.com/egor-denisov/golang-tg-cloud#credits)

## Live Demo
Live demo is available here: 
- [Telegram bot](https://t.me/StorageTest1Bot);
- Web version [currently unavailable].
## Final product
- Startpage lightmode:
![](https://sun9-east.userapi.com/sun9-23/s/v1/ig2/3VWoMe1ZvQJgA0tkF7tpffCN-Gi_kWHBy5JkAgJaBOMH507KWkV87GYnTrRg5_Z0rogZWjuKckvPP9l0fMjgTiDq.jpg?size=1918x930&quality=95&type=album "Start page lightmode")

- Startpage darkmode:
![](https://sun9-north.userapi.com/sun9-80/s/v1/ig2/BasdB0MbfeCsr1KphBKqEqFGHP4z3ar_IsmuIgrtSSfncIkARqar6D-Xl52JsjktJERYcW2Ja0CeJowa-U2xvkaQ.jpg?size=1911x920&quality=95&type=album "Start page darkmode")

- Ships placement page:
![](https://sun9-west.userapi.com/sun9-68/s/v1/ig2/jCUxjO4MKKgyvnHoSCYquzt4esWGZEdtPy4QYKJ4ROlNIE5rz7dyL3FDgiqC3Exc7QF0tX4u3ahTAAfTwpY6mhqn.jpg?size=1914x917&quality=95&type=album "Ships placement page")

- Battle page:
![](https://sun9-east.userapi.com/sun9-25/s/v1/ig2/wneNRVCZIsHxVwyIZEqUxQ8gpdCErEaJ-zUfasAZFAg9LfDTNGeFVboCEOlfmABPI8p3_TeNa_SXJ7Yh4qMFqfWn.jpg?size=1913x923&quality=95&type=album "Battle page")


## Running the project
1) For starting Telegram version you need get token from https://t.me/BotFather. 
2) Further you need to edit all the files named custom.go and set your data there.
3) To create necessary tables, in your PostgreSQL database run commands from sql/query.sql.
4) Finally, in cmd run this command:

```
go build
```
Аfter executing, you can run main.exe file and write /start to your bot.

## About the app

This part is in development...

## Features

This paragraph contains some of the features of this project:

This part is in development...

## Dependencies
- Golang
- PostgreSQL

## Credits
- [tgbotapi](github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [postgres]("github.com/lib/pq")
