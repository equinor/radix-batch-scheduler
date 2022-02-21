package main

import (
    "github.com/equinor/radix-batch-scheduler/models/env"
    log "github.com/sirupsen/logrus"
)

func main() {

    env := env.NewEnv()
    if err := env.ValidateExpected(); err != nil {
        log.Error(err)
        return
    }

    log.Infof("Start the batch %s for the component: %s, deployment: %s", env.BatchName, env.Common.RadixComponentName,
        env.Common.RadixDeploymentName)

    log.Info("Done")
}
