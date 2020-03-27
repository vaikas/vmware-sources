/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mattmoor/vmware-sources/pkg/vsphere"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	apistest "knative.dev/pkg/apis/testing"
)

func TestVSphereBindingDuckTypes(t *testing.T) {
	tests := []struct {
		name string
		t    duck.Implementable
	}{{
		name: "conditions",
		t:    &duckv1.Conditions{},
	}, {
		name: "binding",
		t:    &duckv1alpha1.Binding{},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := duck.VerifyType(&VSphereBinding{}, test.t)
			if err != nil {
				t.Errorf("VerifyType(VSphereBinding, %T) = %v", test.t, err)
			}
		})
	}
}

func TestVSphereBindingGetGroupVersionKind(t *testing.T) {
	r := &VSphereBinding{}
	want := schema.GroupVersionKind{
		Group:   "sources.knative.dev",
		Version: "v1alpha1",
		Kind:    "VSphereBinding",
	}
	if got := r.GetGroupVersionKind(); got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestVSphereBindingUndo(t *testing.T) {
	tests := []struct {
		name string
		in   *duckv1.WithPod
		want *duckv1.WithPod
	}{{
		name: "nothing to remove",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
						}},
					},
				},
			},
		},
	}, {
		name: "lots to remove",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "GOVC_URL",
								Value: "http://localhost:8080",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name:  "GOVC_USERNAME",
								Value: "user",
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "GOVC_URL",
								Value: "http://localhost:8080",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "GOVC_PASSWORD",
								Value: "pass",
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "GOVC_URL",
								Value: "http://localhost:8080",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "GOVC_INSECURE",
								Value: "true",
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Volumes: []corev1.Volume{{
							Name: vsphere.VolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "no-matter",
								},
							},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}},
							VolumeMounts: []corev1.VolumeMount{},
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}},
							VolumeMounts: []corev1.VolumeMount{},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "BAZ",
								Value: "INGA",
							}},
							VolumeMounts: []corev1.VolumeMount{},
						}},
						Volumes: []corev1.Volume{},
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.in
			sb := &VSphereBinding{}
			sb.Undo(context.Background(), got)

			if !cmp.Equal(got, test.want) {
				t.Errorf("Undo (-want, +got): %s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestVSphereBindingDo(t *testing.T) {
	url := apis.URL{
		Scheme: "http",
		Host:   "vmware.com",
	}
	secretName := "ssssshhhh-dont-tell"
	vsb := &VSphereBinding{
		Spec: VSphereBindingSpec{
			VAuthSpec: VAuthSpec{
				Address: url,
				SecretRef: corev1.LocalObjectReference{
					Name: secretName,
				},
			},
		},
	}

	tests := []struct {
		name string
		in   *duckv1.WithPod
		want *duckv1.WithPod
	}{{
		name: "nothing to add",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "GOVC_URL",
								Value: url.String(),
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Volumes: []corev1.Volume{{
							Name: vsphere.VolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretName,
								},
							},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "GOVC_URL",
								Value: url.String(),
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Volumes: []corev1.Volume{{
							Name: vsphere.VolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretName,
								},
							},
						}},
					},
				},
			},
		},
	}, {
		name: "fix the URI",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "GOVC_URL",
								Value: "the wrong value",
							}, {
								Name:  "GOVC_INSECURE",
								Value: `{"extensions":{"wrong":"value"}}`,
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "GOVC_URL",
								Value: url.String(),
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Volumes: []corev1.Volume{{
							Name: vsphere.VolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretName,
								},
							},
						}},
					},
				},
			},
		},
	}, {
		name: "lots to add",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "BAZ",
								Value: "INGA",
							}},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "GOVC_URL",
								Value: url.String(),
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "GOVC_URL",
								Value: url.String(),
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "GOVC_URL",
								Value: url.String(),
							}, {
								Name:  "GOVC_INSECURE",
								Value: "false",
							}, {
								Name: "GOVC_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthUsernameKey,
									},
								},
							}, {
								Name: "GOVC_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secretName,
										},
										Key: corev1.BasicAuthPasswordKey,
									},
								},
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      vsphere.VolumeName,
								ReadOnly:  true,
								MountPath: vsphere.MountPath,
							}},
						}},
						Volumes: []corev1.Volume{{
							Name: vsphere.VolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretName,
								},
							},
						}},
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.in

			ctx := context.Background()

			vsb := vsb.DeepCopy()
			vsb.Do(ctx, got)

			if !cmp.Equal(got, test.want) {
				t.Errorf("Do (-want, +got): %s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestTypicalBindingFlow(t *testing.T) {
	r := &VSphereBindingStatus{}
	r.InitializeConditions()
	apistest.CheckConditionOngoing(r, VSphereBindingConditionReady, t)

	r.MarkBindingUnavailable("Foo", "Bar")
	apistest.CheckConditionFailed(r, VSphereBindingConditionReady, t)

	r.MarkBindingAvailable()
	// After all of that, we're finally ready!
	apistest.CheckConditionSucceeded(r, VSphereBindingConditionReady, t)
}