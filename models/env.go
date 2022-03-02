package models

import (
	"fmt"
	"os"

	"github.com/equinor/radix-common/utils/errors"
	"github.com/equinor/radix-job-scheduler/models"
)

// Env instance variables
type Env struct {
	Common                       *models.Env
	BatchName                    string
	BatchScheduleDescriptionPath string
}

//New Constructor of Env
func New() *Env {
	return &Env{
		Common:                       models.NewEnv(),
		BatchName:                    os.Getenv("RADIX_BATCH_NAME"),
		BatchScheduleDescriptionPath: os.Getenv("RADIX_BATCH_SCHEDULE_DESCRIPTION_PATH"),
	}
}

//ValidateExpected ValidateExpected environment variables
func (env *Env) ValidateExpected() error {
	var errs []error
	if len(env.BatchScheduleDescriptionPath) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_BATCH_SCHEDULE_DESCRIPTION_PATH"))
	}
	if len(env.BatchName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_BATCH_NAME"))
	}
	if len(env.Common.RadixAppName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_APP"))
	}
	if len(env.Common.RadixComponentName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_COMPONENT"))
	}
	if len(env.Common.RadixDefaultCpuLimit) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIXOPERATOR_APP_ENV_LIMITS_DEFAULT_CPU"))
	}
	if len(env.Common.RadixDefaultMemoryLimit) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIXOPERATOR_APP_ENV_LIMITS_DEFAULT_MEMORY"))
	}
	if len(env.Common.RadixDNSZone) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_DNS_ZONE"))
	}
	if len(env.Common.RadixClusterName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_CLUSTERNAME"))
	}
	if len(env.Common.RadixEnvironment) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_ENVIRONMENT"))
	}
	if len(env.Common.RadixDeploymentName) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_DEPLOYMENT"))
	}
	if len(env.Common.RadixContainerRegistry) == 0 {
		errs = append(errs, fmt.Errorf("missed environment variable RADIX_CONTAINER_REGISTRY"))
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Concat(errs)
}
