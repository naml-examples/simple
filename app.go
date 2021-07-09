//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//    ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗
//    ████╗  ██║██╔═████╗██║   ██║██╔══██╗
//    ██╔██╗ ██║██║██╔██║██║   ██║███████║
//    ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║
//    ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║
//    ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝

package app

import (
	"context"
	"fmt"

	naml2 "github.com/kris-nova/naml"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Version string

type MySampleApp struct {
	metav1.ObjectMeta
	exampleString string
	exampleInt    int
	description   string
}

func New(namespace string, name string, description string, exampleString string, exampleInt int) *MySampleApp {
	return &MySampleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			ResourceVersion: "v1.0.0",
			Labels: map[string]string{
				"k8s-app":       "mysampleapp",
				"app":           "mysampleapp",
				"example-label": exampleString,
				"description":   "short-description-of-your-app",
			},
			Annotations: map[string]string{
				"beeps": "boops",
			},
		},
		exampleInt:    exampleInt,
		exampleString: exampleString,
		description:   description,
	}
}

func (v *MySampleApp) Install(client *kubernetes.Clientset) error {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: v.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: naml2.I32p(int32(v.exampleInt)),
			Selector: &metav1.LabelSelector{
				MatchLabels: v.Labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: v.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  v.Name,
							Image: "busybox",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := client.AppsV1().Deployments(v.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to install deployment in Kubernetes: %v", err)
	}
	return nil
}

func (v *MySampleApp) Uninstall(client *kubernetes.Clientset) error {
	return client.AppsV1().Deployments(v.Namespace).Delete(context.TODO(), v.Name, metav1.DeleteOptions{})
}

func (v *MySampleApp) Meta() *metav1.ObjectMeta {
	return &v.ObjectMeta
}

func (v *MySampleApp) Description() string {
	return v.description
}
