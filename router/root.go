package router

import (
    "fmt"
    "regexp"
    "strconv"
    //"encoding/json"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    "MessageForwarder/router/rule"
    "MessageForwarder/router/word"
    "MessageForwarder/models/userdata"
    ruledb "MessageForwarder/models/rule"
    worddb "MessageForwarder/models/word"
)

func Init() (router *commandrouter.Router) {
    router = new(commandrouter.Router)
    router.Set("", root)
    router.Set("getchatid", getchatid)
    router.SetGroup("rule", rule.Init())
    router.SetGroup("word", word.Init())
    router.Set("cancel", cancel)
    return router
}

func root(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    text := update.Message.Text
    //chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    rules, err := ruledb.GetAll(fmt.Sprintf("%d", fromuid))
    if logger.Log(err) {
        return
    }
    //rulesname := []string{}
    for _, nowrule := range rules {
        forward := nowrule.Whitelist == "" && nowrule.Blacklist == ""
        if nowrule.Whitelist != "" {
            nowword, err := worddb.Get(fmt.Sprintf("%d", fromuid), nowrule.Whitelist)
            if logger.Log(err) {
                continue
            }
            re, err := regexp.Compile(nowword.Word)
            if logger.Log(err) {
                continue
            }
            if re.FindString(text) != "" {
                forward = true
            }
        } else if nowrule.Blacklist != "" {
            nowword, err := worddb.Get(fmt.Sprintf("%d", fromuid), nowrule.Blacklist)
            if logger.Log(err) {
                continue
            }
            re, err := regexp.Compile(nowword.Word)
            if logger.Log(err) {
                continue
            }
            if re.FindString(text) == "" {
                forward = true
            }
        }
        if forward {
            for _, replace := range nowrule.Replacelist {
                from, err := worddb.Get(fmt.Sprintf("%d", fromuid), replace.From)
                if logger.Log(err) {
                    continue
                }
                to, err := worddb.Get(fmt.Sprintf("%d", fromuid), replace.To)
                if logger.Log(err) {
                    continue
                }
                re, err := regexp.Compile(from.Word)
                if logger.Log(err) {
                    continue
                }
                text = re.ReplaceAllString(text, to.Word)
            }
            for _, chatidstr := range nowrule.ChatIDlist {
                chatid, err := strconv.ParseInt(chatidstr, 10, 64)
                if logger.Log(err) {
                    continue
                }
                msg := tgbotapi.NewMessage(chatid, text)
                msg.Entities = update.Message.Entities
                _, err = bot.Send(msg)
                logger.Log(err)
            }
        }
    }

    //msgjson, _ := json.Marshal(update.Message)
    //fmt.Println(string(msgjson))
    //msg.ParseMode = tgbotapi.ModeMarkdownV2
    //_, err = bot.Send(msg)
    //logger.Log(err)
}

func getchatid(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    //text := update.Message.Text
    chatid := update.Message.Chat.ID
    reply := fmt.Sprintf("ChatID: `%d`", chatid)
    msg := tgbotapi.NewMessage(chatid, reply)
    msg.ParseMode = tgbotapi.ModeMarkdownV2
    _, err := bot.Send(msg)
    logger.Log(err)
}

func cancel(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    //text := update.Message.Text
    chatid := update.Message.Chat.ID
    fromuid := update.Message.From.ID
    nowuserdata, err := userdata.Get(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
    if logger.Log(err) {
        return
    }
    if nowuserdata == nil {
        _, err = bot.Send(tgbotapi.NewMessage(chatid, "No active command to cancel."))
        logger.Log(err)
    } else {
        _ = userdata.Delete(fmt.Sprintf("%d", fromuid), fmt.Sprintf("%d", chatid))
        msg := tgbotapi.NewMessage(chatid, fmt.Sprintf("The command `%s` has been cancelled\\.", nowuserdata.Data["mode"]))
        msg.ParseMode = tgbotapi.ModeMarkdownV2
        _, err = bot.Send(msg)
        logger.Log(err)
    }
}
