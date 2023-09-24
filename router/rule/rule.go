package rule

import (
    "fmt"
    "strings"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    "MessageForwarder/router/rule/set"
    "MessageForwarder/router/rule/del"
    "MessageForwarder/models/rule"
)

func Init() (router *commandrouter.Router) {
    router = new(commandrouter.Router)
    router.Set("", root)
    router.SetGroup("set", set.Init())
    router.SetGroup("del", del.Init())
    return router
}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    //text := update.Message.Text
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    if !update.Message.Chat.IsPrivate() {
        _, err := bot.Send(tgbotapi.NewMessage(chatid, "Please use this command in private chat."))
        logger.Log(err)
        return
    }
    rules, err := rule.GetAll(fmt.Sprintf("%d", fromuid))
    if logger.Log(err) {
        return
    }
    _, err = bot.Send(tgbotapi.NewMessage(chatid, "Showing list of rules: "))
    logger.Log(err)
    for _, nowrule := range rules {
        reply := fmt.Sprintf("name: `%s`\n" +
                             "to: `%s`\n", nowrule.Name, strings.Join(nowrule.ChatIDlist, "`, `"))
        if nowrule.Whitelist != "" {
            reply += fmt.Sprintf("whitelist: `%s`\n", nowrule.Whitelist)
        } else if nowrule.Blacklist != "" {
            reply += fmt.Sprintf("blacklist: `%s`\n", nowrule.Blacklist)
        }
        if nowrule.Replacelist != nil {
            reply += "replace: \n"
            for _, replace := range nowrule.Replacelist {
                reply += fmt.Sprintf("from: `%s`\n" + 
                                     "to: `%s`\n", replace.From, replace.To)
            }
        }
        msg := tgbotapi.NewMessage(chatid, reply)
        msg.ParseMode = tgbotapi.ModeMarkdownV2
        _, err = bot.Send(msg)
        logger.Log(err)
    }
    msg := tgbotapi.NewMessage(chatid, "What do you want to do?")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Set rule", "rule.set"),
            tgbotapi.NewInlineKeyboardButtonData("Delete rule", "rule.del"),
        ),
    )
    _, err = bot.Send(msg)
}

