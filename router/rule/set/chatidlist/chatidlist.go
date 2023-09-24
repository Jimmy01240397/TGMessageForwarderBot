package chatidlist

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    //"MessageForwarder/models/rule"
    "MessageForwarder/models/userdata"
)

func Init() (router *commandrouter.Router) {
    router = new(commandrouter.Router)
    router.Set("", root)
    router.Set("finish", finish)
    return router
}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    text := update.Message.Text
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        return
    }
    _, err = bot.Request(tgbotapi.NewDeleteMessage(chatid, int(nowuserdata.Data["msgid"].(float64))))
    logger.Log(err)
    msg := tgbotapi.NewMessage(chatid, "Enter chatid that you want to forward message to or press finish to finish chatidlist.")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Finish", "rule.set.chatidlist.finish"),
        ),
    )
    msgdata, err := bot.Send(msg)
    logger.Log(err)
    nowuserdata.Data["mode"] = "rule.set.chatidlist"
    nowuserdata.Data["msgid"] = msgdata.MessageID
    nowuserdata.Data["chatidlist"] = append(nowuserdata.Data["chatidlist"].([]any), text)
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), nowuserdata.Data)
    if logger.Log(err) {
        return
    }
}

func finish(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.CallbackQuery.Message.Chat.ID
    msgid := update.CallbackQuery.Message.MessageID
    fromuid := update.CallbackQuery.From.ID
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
        return
    }
    nowuserdata.Data["mode"] = "rule.set.nothing"
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), nowuserdata.Data)
    if logger.Log(err) {
        _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
        return
    }
    msgedit := tgbotapi.NewEditMessageTextAndMarkup(chatid, msgid, "Do you want to set whitelist or blacklist?", tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Whitelist", "rule.set.word.white"),
            tgbotapi.NewInlineKeyboardButtonData("Blacklist", "rule.set.word.black"),
            tgbotapi.NewInlineKeyboardButtonData("No", "rule.set.word.end"),
        ),
    ))
    _, err = bot.Send(msgedit)
    logger.Log(err)
}
