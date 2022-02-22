package main

import (
	"github.com/equinor/radix-batch-scheduler/models"
	log "github.com/sirupsen/logrus"
)

func main() {

	env := models.New()
	if err := env.ValidateExpected(); err != nil {
		log.Error(err.Error())
		return
	}

	log.Infof("Start the batch '%s' for the component: '%s', deployment: '%s'", env.BatchName, env.Common.RadixComponentName,
		env.Common.RadixDeploymentName)

	//kubeClient, radixClient, _, secretProviderClient := utils.GetKubernetesClient()
	//kubeUtil, _ := kube.New(kubeClient, radixClient, secretProviderClient)
	//job := jobApi.New(env.Common, kubeUtil, kubeClient, radixClient)
	//
	//log.Debugf("Run batch %s", job.)

	log.Info("Done")
}
