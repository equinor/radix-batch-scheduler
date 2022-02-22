package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/equinor/radix-batch-scheduler/models"
	jobApi "github.com/equinor/radix-job-scheduler/api/jobs"
	schedulerDefaults "github.com/equinor/radix-job-scheduler/defaults"
	apiModels "github.com/equinor/radix-job-scheduler/models"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	"github.com/equinor/radix-operator/pkg/apis/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	env := models.New()
	if err := env.ValidateExpected(); err != nil {
		log.Error(err.Error())
		return
	}

	log.Infof("Requested a batch '%s' for the component: '%s', deployment: '%s'", env.BatchName,
		env.Common.RadixComponentName, env.Common.RadixDeploymentName)

	batchScheduleDescription, err := getBatchScheduleDescription(env)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(batchScheduleDescription.JobScheduleDescriptions) == 0 {
		log.Info("JobScheduleDescriptions list is empty")
		return
	}

	err = runJobs(env, batchScheduleDescription)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("All jobs processed.")
}

func getBatchScheduleDescription(env *models.Env) (*apiModels.BatchScheduleDescription, error) {
	secretName := schedulerDefaults.GetBatchScheduleDescriptionSecretName(env.BatchName)
	secretPath := path.Join(schedulerDefaults.BatchSecretsMountPath, secretName)
	if _, err := os.Stat(secretPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("missing secret file %s, possible missed the secret %s", secretPath, secretName)
	}
	batchScheduleDescriptionBuff, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, err
	}
	batchScheduleDescription := &apiModels.BatchScheduleDescription{}
	err = json.Unmarshal(batchScheduleDescriptionBuff, batchScheduleDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to read the BatchScheduleDescription from the secret: %v", err)
	}
	return batchScheduleDescription, nil
}

func runJobs(env *models.Env, batchScheduleDescription *apiModels.BatchScheduleDescription) error {
	kubeClient, radixClient, _, secretProviderClient := utils.GetKubernetesClient()
	kubeUtil, err := kube.New(kubeClient, radixClient, secretProviderClient)
	if err != nil {
		return err
	}
	jobModel := jobApi.New(env.Common, kubeUtil, kubeClient, radixClient)
	log.Info("Run the batch '%s'", env.BatchName)
	var jobStatuses []*apiModels.JobStatus
	for _, jobScheduleDescription := range batchScheduleDescription.JobScheduleDescriptions {
		jobStatuses = append(jobStatuses, runJob(jobModel, jobScheduleDescription))
	}
	return nil
}

func runJob(jobModel jobApi.Job, jobScheduleDescription apiModels.JobScheduleDescription) *apiModels.JobStatus {
	log.Infof("Start the job '%s'", jobScheduleDescription.JobId)
	jobStatus, err := jobModel.CreateJob(&jobScheduleDescription)
	if err != nil {
		log.Errorf("failed start the job '%s': %v", jobScheduleDescription.JobId, err)
		return &apiModels.JobStatus{
			JobId:   jobScheduleDescription.JobId,
			Status:  "failed",
			Message: err.Error(),
		}
	}
	log.Info("job '%s' has been started, job name: %s", jobScheduleDescription.JobId, jobStatus.Name)
	return jobStatus
}
