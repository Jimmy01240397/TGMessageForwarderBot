package word

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
    router.Set("white", whitelist)
    router.Set("black", blacklist)
    router.Set("set", setlist)
    router.Set("end", wordend)
    return router
}

func whitelist(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    list("whitelist", bot, update)
}

func blacklist(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    list("blacklist", bot, update)
}

func list(listname string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.CallbackQuery.Message.Chat.ID
    msgid := update.CallbackQuery.Message.MessageID
    fromuid := update.CallbackQuery.From.ID
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
        return
    }
    nowuserdata.Data["mode"] = "rule.set.word.set"
    nowuserdata.Data["listtype"] = listname
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), nowuserdata.Data)
    if logger.Log(err) {
        _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
        return
    }
    msgedit := tgbotapi.NewEditMessageText(chatid, msgid, fmt.Sprintf("Enter your word name for %s", listname))
    _, err = bot.Send(msgedit)
    logger.Log(err)
}

func setlist(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    text := update.Message.Text
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        return
    }
    nowuserdata.Data["listname"] = text
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), nowuserdata.Data)
    if logger.Log(err) {
        return
    }
    wordend(bot, update)
}

func wordend(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var chatid int64
    var fromuid int64
    var msgid int
    if update.Message != nil {
        chatid = update.Message.Chat.ID
        fromuid = update.Message.From.ID
    } else if update.CallbackQuery != nil {
        chatid = update.CallbackQuery.Message.Chat.ID
        msgid = update.CallbackQuery.Message.MessageID
        fromuid = update.CallbackQuery.From.ID
    }
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        if update.CallbackQuery != nil {
            _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
        }
        return
    }
    nowuserdata.Data["mode"] = "rule.set.nothing"
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), nowuserdata.Data)
    if logger.Log(err) {
        if update.CallbackQuery != nil {
            _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
        }
        return
    }
    replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Yes", "rule.set.replace"),
            tgbotapi.NewInlineKeyboardButtonData("No", "rule.set.replace.finish"),
        ),
    )
    var msg tgbotapi.Chattable
    if update.Message != nil {
        tmp := tgbotapi.NewMessage(chatid, "Do you want to set replace rule?")
        tmp.ReplyMarkup = replyMarkup
        msg = tmp
    } else if update.CallbackQuery != nil {
        msg = tgbotapi.NewEditMessageTextAndMarkup(chatid, msgid, "Do you want to set replace rule?", replyMarkup)
    }
    _, err = bot.Send(msg)
    logger.Log(err)
}
