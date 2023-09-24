package commandrouter

import (
    "fmt"
    "strings"

    //"golang.org/x/exp/slices"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type routedata interface {
    run(fullpath, path string, bot *tgbotapi.BotAPI, update tgbotapi.Update)
}

type Router struct {
    data map[string]routedata
}

type routefunc func(*tgbotapi.BotAPI, tgbotapi.Update)

type routefuncdata struct {
    runner routefunc
}

func (router *Router) Set(path string, route routefunc) {
    if router.data == nil {
        router.data = make(map[string]routedata)
    }
    router.data[path] = &routefuncdata{
        runner: route,
    }
}

func (router *Router) SetGroup(path string, routegroup *Router) {
    router.data[path] = routegroup
}

func (router *Router) Run(path string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    router.run(path, path, bot, update)
}

func (router *Router) run(fullpath, path string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    pathslice := strings.Split(path, ".")
    for i := len(pathslice); i >= 0; i-- {
        nowroute, exist := router.data[strings.Join(pathslice[:i], ".")]
        if exist {
            nowroute.run(fullpath, strings.Join(pathslice[i:], "."), bot, update)
            break
        }
    }
}

func (exec *routefuncdata) run(fullpath, path string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    if path != "" {
        var chatid int64
        if update.Message != nil {
            chatid = update.Message.Chat.ID
        } else if update.CallbackQuery != nil {
            chatid = update.CallbackQuery.Message.Chat.ID
        } else {
            return
        }
        msg := tgbotapi.NewMessage(chatid, fmt.Sprintf("Command `%s` not found\\!", fullpath))
        msg.ParseMode = tgbotapi.ModeMarkdownV2
        _, _ = bot.Send(msg)
        return
    }
    exec.runner(bot, update)
}
