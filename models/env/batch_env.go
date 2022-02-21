package env

import (
    "fmt"
    "os"

    commonErrors "github.com/equinor/radix-common/utils/errors"
    schedulerModels "github.com/equinor/radix-job-scheduler/models"
)

type Env struct {
    Common    *schedulerModels.Env
    BatchName string
}

//NewEnv CConstructor of Env
func NewEnv() *Env {
    return &Env{
        Common:    schedulerModels.NewEnv(),
        BatchName: os.Getenv("RADIX_BATCH_NAME"),
    }
}

//ValidateExpected ValidateExpected environment variables
func (env *Env) ValidateExpected() error {
    var errs []error
    if len(env.BatchName) == 0 {
        errs = append(errs, fmt.Errorf("missed environment variable RADIX_BATCH_JOB"))
    }
    if len(env.Common.RadixAppName) == 0 {
        errs = append(errs, fmt.Errorf("missed environment variable RADIX_APP"))
    }
    if len(env.Common.RadixComponentName) == 0 {
        errs = append(errs, fmt.Errorf("missed environment variable RADIX_COMPONENT"))
    }
    return commonErrors.Concat(errs)
}
