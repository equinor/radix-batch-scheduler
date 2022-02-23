package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/equinor/radix-batch-scheduler/models"
	jobApi "github.com/equinor/radix-job-scheduler/api/jobs"
	apiModels "github.com/equinor/radix-job-scheduler/models"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	"github.com/equinor/radix-operator/pkg/apis/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
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
	batchScheduleDescriptionPath := env.BatchScheduleDescriptionPath
	if _, err := os.Stat(batchScheduleDescriptionPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("missing secret file %s, possible missed the secret for the BatchScheduleDescription",
			batchScheduleDescriptionPath)
	}
	batchScheduleDescriptionBuff, err := ioutil.ReadFile(batchScheduleDescriptionPath)
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
	log.Infof("Run the batch '%s' of %d jobs", env.BatchName, len(batchScheduleDescription.JobScheduleDescriptions))
	for _, jobScheduleDescription := range batchScheduleDescription.JobScheduleDescriptions {
		runJob(jobModel, jobScheduleDescription)
	}
	return nil
}

func runJob(jobModel jobApi.Job, jobScheduleDescription apiModels.JobScheduleDescription) {
	log.Infof("Start the job '%s'", jobScheduleDescription.JobId)
	jobStatus, err := jobModel.CreateJob(&jobScheduleDescription)
	if err != nil {
		log.Errorf("failed start the job '%s': %v", jobScheduleDescription.JobId, err)
		return
	}
	log.Infof("job '%s' has been started, job name: '%s'", jobScheduleDescription.JobId, jobStatus.Name)
}
