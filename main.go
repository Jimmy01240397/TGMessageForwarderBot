package main

import (
    "fmt"
    "log"
    
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/config"
    //"MessageForwarder/utils/logger"
    "MessageForwarder/models/userdata"
    "MessageForwarder/router"
)

var botrouter *commandrouter.Router

func main() {
    bot, err := tgbotapi.NewBotAPI(config.Token)
    if err != nil {
        panic(err)
    }

    _, _ = bot.Request(tgbotapi.NewSetMyCommands(
        tgbotapi.BotCommand{
            Command: "/getchatid",
            Description: "Get chat id in this room",
        },
        tgbotapi.BotCommand{
            Command: "/rule",
            Description: "Setup forward rule",
        },
        tgbotapi.BotCommand{
            Command: "/word",
            Description: "Setup word data for forward rule",
        },
        tgbotapi.BotCommand{
            Command: "/cancel",
            Description: "Cancel the current operation",
        },
    ))

    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = 60
    updates := bot.GetUpdatesChan(updateConfig)

    botrouter = router.Init()

    for update := range updates {
        go handleUpdate(bot, update)
    }
}

func panicHandler() {
    if err := recover(); err != nil {
        log.Println("panic occurred:", err)
    }
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    defer panicHandler()
    if update.Message != nil {
        if update.Message.IsCommand() {
            botrouter.Run(update.Message.Command(), bot, update)
        } else {
            chatid := update.Message.Chat.ID
            fromuid := update.Message.From.ID
            nowuserdata, _ := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
            if nowuserdata != nil {
                botrouter.Run(nowuserdata.Data["mode"].(string), bot, update)
            } else {
                botrouter.Run("", bot, update)
            }
        }
    }
    if update.CallbackQuery != nil {
        botrouter.Run(update.CallbackQuery.Data, bot, update)
    }
}

