# MessageForwarder

You can use this bot to send anonymous message to your group by send message to this bot.

## Build
1. install golang 1.19
2. Clone this repo and cd into TGMessageForwarderBot
``` bash
git clone https://github.com/Jimmy01240397/TGMessageForwarderBot
cd TGMessageForwarderBot
```
3. Run make
``` bash
make
```

## Run
1. After [Build](#build) cd into `bin` dir
``` bash
cd bin
```
2. Copy .env.sample to .env and write your config
``` bash
cp .env.sample .env
vim .env
```

```
TGBOTTOKEN=<token>
DBNAME=data.db
```
3. Run `tgmsgforwarderbot`
``` bash
./tgmsgforwarderbot
```
