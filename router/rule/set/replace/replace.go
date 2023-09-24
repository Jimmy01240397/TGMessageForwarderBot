package replace

import (
    "fmt"
    "strings"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    "MessageForwarder/models/rule"
    "MessageForwarder/models/userdata"
)

func Init() (router *commandrouter.Router) {
    router = new(commandrouter.Router)
    router.Set("", root)
    router.Set("finish", finish)
    return router
}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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
    nowuserdata, _ := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if nowuserdata != nil && update.Message != nil {
        msgid = int(nowuserdata.Data["msgid"].(float64))
    }
    _, err := bot.Request(tgbotapi.NewDeleteMessage(chatid, msgid))
    logger.Log(err)
    if nowuserdata == nil {
        _ = userdata.Delete(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
        _, err = bot.Send(tgbotapi.NewMessage(chatid, "Something wrong!"))
        logger.Log(err)
        return
    }
    msg := tgbotapi.NewMessage(chatid, "Enter replace rule with format `from_word_data_name,to_word_data_name` or press finish to finish replace rule\\.")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Finish", "rule.set.replace.finish"),
        ),
    )
    msg.ParseMode = tgbotapi.ModeMarkdownV2
    msgdata, err := bot.Send(msg)
    logger.Log(err)
    nowuserdata.Data["mode"] = "rule.set.replace"
    nowuserdata.Data["msgid"] = msgdata.MessageID
    replacelist, exist := nowuserdata.Data["replacelist"]
    if !exist {
        replacelist = []any{}
    }
    if update.Message != nil {
        text := update.Message.Text
        nowuserdata.Data["replacelist"] = append(replacelist.([]any), text)
    }
    err = userdata.Set(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid), nowuserdata.Data)
    if logger.Log(err) {
        return
    }
}

func finish(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatid := update.CallbackQuery.Message.Chat.ID
    msgid := update.CallbackQuery.Message.MessageID
    fromuid := update.CallbackQuery.From.ID
    nowuserdata, _ := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    _ = userdata.Delete(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    var msgedit tgbotapi.EditMessageTextConfig
    if nowuserdata != nil {
        chatidlist := []string{}
        for _, chatid := range nowuserdata.Data["chatidlist"].([]any) {
            chatidlist = append(chatidlist, chatid.(string))
        }
        listtype, exist := nowuserdata.Data["listtype"]
        whitelist := ""
        blacklist := ""
        if exist && listtype.(string) == "whitelist" {
            whitelist = nowuserdata.Data["listname"].(string)
        } else if exist && listtype.(string) == "blacklist" {
            blacklist = nowuserdata.Data["listname"].(string)
        }
        replacelist := []rule.Replace{}
        replacelistdata, exist := nowuserdata.Data["replacelist"]
        if exist {
            for _, replace := range replacelistdata.([]any) {
                replacedata := strings.Split(replace.(string), ",")
                if len(replacedata) == 2 {
                    replacelist = append(replacelist, rule.Replace{
                        From: replacedata[0],
                        To: replacedata[1],
                    })
                }
            }
        }
        err := rule.Set(fmt.Sprintf("%d", fromuid), nowuserdata.Data["name"].(string), chatidlist, whitelist, blacklist, replacelist)
        if logger.Log(err) {
            _, _ = bot.Send(tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!"))
            return
        }
        msgedit = tgbotapi.NewEditMessageText(chatid, msgid, "Rule set success!")
    } else {
        msgedit = tgbotapi.NewEditMessageText(chatid, msgid, "Something wrong!")
    }
    _, _ = bot.Send(msgedit)
}
