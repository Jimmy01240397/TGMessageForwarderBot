package set

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
    router.Set("word", setword)
    return router
}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.CallbackQuery.Message.Chat.ID
    msgid := update.CallbackQuery.Message.MessageID
    fromuid := update.CallbackQuery.From.ID
    err := userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), map[string]any{
        "mode": "word.set.name",
    })
    if logger.Log(err) {
        return
    }
    msgedit := tgbotapi.NewEditMessageText(chatid, msgid, "Enter your word data name that we can identify your word data in your rule.")
    _, err = bot.Send(msgedit)
    logger.Log(err)
}

func name(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    text := update.Message.Text
    err := userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), map[string]any{
        "mode": "word.set.word",
        "name": text,
    })
    if logger.Log(err) {
        return
    }
    msg := tgbotapi.NewMessage(chatid, "Enter a regex for your word data.")
    _, err = bot.Send(msg)
    logger.Log(err)
}

func setword(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    text := update.Message.Text
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        return
    }
    _ = userdata.Delete(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    err = word.Set(fmt.Sprintf("%d", fromuid), nowuserdata.Data["name"].(string), text)
    if logger.Log(err) {
        return
    }
    msg := tgbotapi.NewMessage(chatid, "Word data set success!")
    _, err = bot.Send(msg)
    logger.Log(err)
}
