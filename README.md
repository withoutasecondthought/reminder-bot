# Reminder bot
## Description
small bot to remind you and your friends about messages in chat

## Usage
### Preparation
create bot with **@BotFather** and get **TOKEN**
### Important
go to **bot settings** and disable **privacy mode**,
or you will need to give bot **admin rights**
### For comfort usage
go to **Edit bot** -> **Edit commands** and add commands:
```
help - get more information about bot
remind - <time> <@nickname> <@nickname> ... <@nickname>
```

### Installation
```shell
git clone https://github.com/withoutasecondthought/reminder-bot.git
```

### Setup you env

to start with default settings:
add **TOKEN** to **env**

if you want to **customize** answers you can change **answers** in **env**

### Start

```shell
cd reminder-bot
docker build -t reminder_bot .
docker run -d --name reminder_bot reminder_bot
```
