package api

import (
	"fmt"
	"github.com/equinor/radix-batch-scheduler/models"
	batchApi "github.com/equinor/radix-job-scheduler/api/batches"
	jobApi "github.com/equinor/radix-job-scheduler/api/jobs"
	apiModels "github.com/equinor/radix-job-scheduler/models"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	notCompletedJobStatuses = map[string]bool{"Running": true, "Waiting": true, "Stopping": true}
)

//RunBatchJobs Run batch jobs
func RunBatchJobs(kubeUtil *kube.Kube, env *models.Env, batchScheduleDescription *apiModels.BatchScheduleDescription) error {
	jobHandler := jobApi.New(env.Common, kubeUtil)
	log.Infof("Run the batch '%s' of %d jobs", env.BatchName, len(batchScheduleDescription.JobScheduleDescriptions))
	jobCount := 0
	for _, jobScheduleDescription := range batchScheduleDescription.JobScheduleDescriptions {
		description := jobScheduleDescription
		applyDefaultJobDescriptionProperties(&description, batchScheduleDescription)
		jobCount++
		runJob(env.BatchName, jobHandler, description, jobCount)
	}
	return nil
}

//CheckIfAllJobsAreCompleted Wait while all job complete or failed
func CheckIfAllJobsAreCompleted(batch batchApi.BatchHandler, env *models.Env, done chan bool) {
	log.Debug("Check if all jobs have been completed")
	batchStatus, err := batch.GetBatch(env.BatchName)
	if err != nil {
		log.Error(err)
		done <- true
		return
	}
	for _, jobStatus := range batchStatus.JobStatuses {
		if _, exists := notCompletedJobStatuses[jobStatus.Status]; exists {
			log.Debug("More not-completed jobs exists. Waiting...")
			return
		}
	}
	done <- true
	log.Info("All jobs completed")
}

func applyDefaultJobDescriptionProperties(jobDescription *apiModels.JobScheduleDescription, batchScheduleDescription *apiModels.BatchScheduleDescription) {
	if batchScheduleDescription == nil || batchScheduleDescription.DefaultRadixJobComponentConfig == nil {
		return
	}
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

func runJob(batchName string, jobHandler jobApi.JobHandler, jobScheduleDescription apiModels.JobScheduleDescription,
	jobCount int) {
	jobName := fmt.Sprintf("#%d", jobCount)
	jobId := strings.TrimSpace(jobScheduleDescription.JobId)
	if len(jobId) > 0 {
		jobName = fmt.Sprintf("%s job-id: '%s'", jobName, jobScheduleDescription.JobId)
	}
	log.Infof("Start the job %s", jobName)
	jobStatus, err := jobHandler.CreateJob(&jobScheduleDescription, batchName)
	if err != nil {
		log.Errorf("failed start the job %s: %v", jobName, err)
		return
	}
	log.Infof("job %s has been started, job name: '%s'", jobName, jobStatus.Name)
}
