// Copyright The OpenTelemetry Authors
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

package registry_test

import (
	"context"
	"os"
	"testing"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	ctx    context.Context
	cancel context.CancelFunc

	k8sClient  client.Client
	testScheme *runtime.Scheme = scheme.Scheme
)

func TestMain(m *testing.M) {
	ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()
	utilruntime.Must(registryv1alpha1.AddToScheme(testScheme))

	k8sClient = fake.NewClientBuilder().WithScheme(testScheme).WithObjects(
	// TODO: add objects if needed
	).Build()

	os.Exit(m.Run())
}
