package word

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "MessageForwarder/utils/commandrouter"
    "MessageForwarder/utils/logger"
    "MessageForwarder/router/word/set"
    "MessageForwarder/router/word/del"
    "MessageForwarder/models/word"
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
    words, err := word.GetAll(fmt.Sprintf("%d", fromuid))
    if logger.Log(err) {
        return
    }
    _, err = bot.Send(tgbotapi.NewMessage(chatid, "Showing list of words: "))
    logger.Log(err)
    for _, nowword := range words {
        reply := fmt.Sprintf("name: `%s`\n" +
                             "word: `%s`\n", nowword.Name, nowword.Word)
        msg := tgbotapi.NewMessage(chatid, reply)
        msg.ParseMode = tgbotapi.ModeMarkdownV2
        _, err = bot.Send(msg)
        logger.Log(err)
    }
    msg := tgbotapi.NewMessage(chatid, "What do you want to do?")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Set word", "word.set"),
            tgbotapi.NewInlineKeyboardButtonData("Delete word", "word.del"),
        ),
    )
    _, err = bot.Send(msg)
    logger.Log(err)
}

