package logger

import (
    "log"
    "MessageForwarder/utils/config"
)

func Log(err error) (bool) {
    if err != nil && config.Debug {
        log.Println(err)
    }
    return err != nil
}
