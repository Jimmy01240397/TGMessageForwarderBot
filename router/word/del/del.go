package del

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    //"MessageForwarder/models/rule"
    "MessageForwarder/models/userdata"
    "MessageForwarder/models/word"
)

func Init() (router *commandrouter.Router) {
    router = new(commandrouter.Router)
    router.Set("", root)
    router.Set("name", name)
    //router.SetGroup("set", set.Init())
    //router.SetGroup("del", del.Init())
    return router
}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.CallbackQuery.Message.Chat.ID
    msgid := update.CallbackQuery.Message.MessageID
    fromuid := update.CallbackQuery.From.ID
    err := userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), map[string]any{
        "mode": "word.del.name",
    })
    if logger.Log(err) {
        return
    }
    msgedit := tgbotapi.NewEditMessageText(chatid, msgid, "Enter your word data name that you want to delete.")
    _, err = bot.Send(msgedit)
    logger.Log(err)
}

func name(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    text := update.Message.Text
    _ = userdata.Delete(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    err := word.Delete(fmt.Sprintf("%d", fromuid), text)
    if err != nil {
        msg := tgbotapi.NewMessage(chatid, err.Error())
        _, err = bot.Send(msg)
        logger.Log(err)
        return
    }
    msg := tgbotapi.NewMessage(chatid, "Word data delete success!")
    _, err = bot.Send(msg)
    logger.Log(err)
}
