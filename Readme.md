# MessageForwarder

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
1. Copy .env.sample to .env and write your config
``` bash
cp .env.sample .env
vim .env
```

```
TGBOTTOKEN=<token>
DBNAME=data.db
```
2. Run tgmsgforwarderbot
``` bash
cd bin
./tgmsgforwarderbot
```
