package api

import (
	"fmt"
	"github.com/equinor/radix-batch-scheduler/models"
	jobApi "github.com/equinor/radix-job-scheduler/api/jobs"
	apiModels "github.com/equinor/radix-job-scheduler/models"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	log "github.com/sirupsen/logrus"
	"strings"
)

//RunBatchJobs Run batch jobs
func RunBatchJobs(kubeUtil *kube.Kube, env *models.Env, batchScheduleDescription *apiModels.BatchScheduleDescription) error {
	jobModel := jobApi.New(env.Common, kubeUtil)
	log.Infof("Run the batch '%s' of %d jobs", env.BatchName, len(batchScheduleDescription.JobScheduleDescriptions))
	jobCount := 0
	for _, jobScheduleDescription := range batchScheduleDescription.JobScheduleDescriptions {
		description := jobScheduleDescription
		applyDefaultJobDescriptionProperties(&description, batchScheduleDescription)
		jobCount++
		runJob(jobModel, description, jobCount)
	}
	return nil
}

func applyDefaultJobDescriptionProperties(jobDescription *apiModels.JobScheduleDescription, batchScheduleDescription *apiModels.BatchScheduleDescription) {
	batchDescription := *batchScheduleDescription
	if jobDescription.RadixJobComponentConfig.Node == nil {
		jobDescription.RadixJobComponentConfig.Node = batchDescription.DefaultRadixJobComponentConfig.Node
	}
	if jobDescription.RadixJobComponentConfig.Resources == nil {
		jobDescription.RadixJobComponentConfig.Resources = batchDescription.DefaultRadixJobComponentConfig.Resources
	}
	if jobDescription.RadixJobComponentConfig.TimeLimitSeconds == nil {
		jobDescription.RadixJobComponentConfig.TimeLimitSeconds = batchDescription.DefaultRadixJobComponentConfig.TimeLimitSeconds
	}
}

func runJob(jobModel jobApi.Job, jobScheduleDescription apiModels.JobScheduleDescription, jobCount int) {
	jobName := fmt.Sprintf("#%d", jobCount)
	jobId := strings.TrimSpace(jobScheduleDescription.JobId)
	if len(jobId) > 0 {
		jobName = fmt.Sprintf("%s job-id: '%s'", jobName, jobScheduleDescription.JobId)
	}
	log.Infof("Start the job %s", jobName)
	jobStatus, err := jobModel.CreateJob(&jobScheduleDescription)
	if err != nil {
		log.Errorf("failed start the job %s: %v", jobName, err)
		return
	}
	log.Infof("job %s has been started, job name: '%s'", jobName, jobStatus.Name)
}
