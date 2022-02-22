package models

import (
	"fmt"
	"os"

	"github.com/equinor/radix-common/utils/errors"
	"github.com/equinor/radix-job-scheduler/models"
)

// Env instance variables
type Env struct {
	Common    *models.Env
	BatchName string
}

//New Constructor of Env
func New() *Env {
	return &Env{
		Common:    models.NewEnv(),
		BatchName: os.Getenv("RADIX_BATCH_NAME"),
	}
}

//ValidateExpected ValidateExpected environment variables
func (env *Env) ValidateExpected() error {
	var errs []error
	if len(env.BatchName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_BATCH_NAME"))
	}
	//if len(env.Common.RadixAppName) == 0 {
	//    errs = append(errs, fmt.Errorf("missed environment variable RADIX_APP"))
	//}
	if len(env.Common.RadixComponentName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_COMPONENT"))
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Concat(errs)
}
