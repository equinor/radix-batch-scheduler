package main

import (
	"encoding/json"
	"errors"
	"fmt"
	batchApi "github.com/equinor/radix-job-scheduler/api/batches"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	"github.com/equinor/radix-operator/pkg/apis/utils"
	"github.com/go-co-op/gocron"
	"io/ioutil"
	"os"
	"time"

	"github.com/equinor/radix-batch-scheduler/api"
	"github.com/equinor/radix-batch-scheduler/models"
	apiModels "github.com/equinor/radix-job-scheduler/models"
	log "github.com/sirupsen/logrus"
)

func main() {
	env := models.New()
	log.SetLevel(log.DebugLevel)

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
		log.Info("required JobScheduleDescriptions list is empty")
		return
	}

	kubeUtil := getKubeUtil()

	err = api.RunBatchJobs(kubeUtil, env, batchScheduleDescription)
	if err != nil {
		log.Error(err.Error())
		return
	}

	waitWhileAllBatchJobsAreCompleted(kubeUtil, env)

	log.Info("All jobs processed.")
}

func waitWhileAllBatchJobsAreCompleted(kubeUtil *kube.Kube, env *models.Env) {
	batch := batchApi.New(env.Common, kubeUtil, kubeUtil.KubeClient(), kubeUtil.RadixClient())
	done := make(chan bool)
	waitScheduler := gocron.NewScheduler(time.UTC).SingletonMode()
	waitScheduler.Every(10).Seconds().Do(func() {
		api.CheckIfAllJobsAreCompleted(batch, env, done)
	})
	waitScheduler.StartAsync()
	<-done
	waitScheduler.Stop()
}

func getKubeUtil() *kube.Kube {
	kubeClient, radixClient, _, secretProviderClient := utils.GetKubernetesClient()
	kubeUtil, _ := kube.New(kubeClient, radixClient, secretProviderClient)
	return kubeUtil
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
