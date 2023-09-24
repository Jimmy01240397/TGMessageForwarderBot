package set

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    //"MessageForwarder/models/rule"
    "MessageForwarder/models/userdata"
    "MessageForwarder/router/rule/set/chatidlist"
    "MessageForwarder/router/rule/set/word"
    "MessageForwarder/router/rule/set/replace"
)

func Init() (router *commandrouter.Router) {
    router = new(commandrouter.Router)
    router.Set("", root)
    router.Set("nothing", nothing)
    router.Set("name", name)
    router.SetGroup("chatidlist", chatidlist.Init())
    router.SetGroup("word", word.Init())
    router.SetGroup("replace", replace.Init())
    return router
}

func nothing(bot *tgbotapi.BotAPI, update tgbotapi.Update) {}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.CallbackQuery.Message.Chat.ID
    msgid := update.CallbackQuery.Message.MessageID
    fromuid := update.CallbackQuery.From.ID
    err := userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), map[string]any{
        "mode": "rule.set.name",
    })
    if logger.Log(err) {
        return
    }
    msgedit := tgbotapi.NewEditMessageText(chatid, msgid, "Enter your rule name.")
    _, err = bot.Send(msgedit)
    logger.Log(err)
}

func name(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    text := update.Message.Text
    msg := tgbotapi.NewMessage(chatid, "Enter chatid that you want to forward message to or press finish to finish chatidlist.")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Finish", "rule.set.chatidlist.finish"),
        ),
    )
    msgdata, err := bot.Send(msg)
    logger.Log(err)
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), map[string]any{
        "mode": "rule.set.chatidlist",
        "name": text,
        "msgid": msgdata.MessageID,
        "chatidlist": []any{},
    })
    if logger.Log(err) {
        return
    }
}

