package cmd

import (
	"fmt"

	jenkinsv1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx/pkg/client/clientset/versioned"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/pkg/errors"
)

func (o *CommonOptions) registerEnvironmentCRD() error {
	apisClient, err := o.ApiExtensionsClient()
	if err != nil {
		return err
	}
	err = kube.RegisterEnvironmentCRD(apisClient)
	return err
}

// modifyDevEnvironment performs some mutation on the Development environment to modify team settings
func (o *CommonOptions) modifyDevEnvironment(jxClient versioned.Interface, ns string,
	fn func(env *jenkinsv1.Environment) error) error {
	env, err := kube.EnsureDevEnvironmentSetup(jxClient, ns)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure that dev environment is setup for namespace '%s'", ns)
	}
	if env == nil {
		return fmt.Errorf("No Development environment found in namespace %s", ns)
	}
	err = fn(env)
	if err != nil {
		return errors.Wrap(err, "failed to call the callback function for dev environment")
	}
	_, err = jxClient.JenkinsV1().Environments(ns).PatchUpdate(env)
	if err != nil {
		return fmt.Errorf("Failed to update Development environment in namespace %s: %s", ns, err)
	}
	return nil
}
