/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package builder

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/apache/camel-k/v2/pkg/util/io"
	"github.com/apache/camel-k/v2/pkg/util/kubernetes"
	"github.com/apache/camel-k/v2/pkg/util/log"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	ctrl "sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
	"github.com/apache/camel-k/v2/pkg/util"
	"github.com/apache/camel-k/v2/pkg/util/defaults"
)

const (
	// ContextDir is the directory used to package the container.
	ContextDir = "context"
	// DeploymentDir is the directory used in the runtime application to deploy the artifacts.
	DeploymentDir = "/deployments"
	// DependenciesDir is the directory used to store required dependencies.
	DependenciesDir = "dependencies"
)

func init() {
	registerSteps(Image)
}

type imageSteps struct {
	IncrementalImageContext Step
	NativeImageContext      Step
	StandardImageContext    Step
	ExecutableDockerfile    Step
	JvmDockerfile           Step
}

// Image used to export the steps available on an Image building process.
var Image = imageSteps{
	IncrementalImageContext: NewStep(ApplicationPackagePhase, incrementalImageContext),
	NativeImageContext:      NewStep(ApplicationPackagePhase, nativeImageContext),
	StandardImageContext:    NewStep(ApplicationPackagePhase, standardImageContext),
	ExecutableDockerfile:    NewStep(ApplicationPackagePhase+1, executableDockerfile),
	JvmDockerfile:           NewStep(ApplicationPackagePhase+1, jvmDockerfile),
}

type artifactsSelector func(ctx *builderContext) error

func nativeImageContext(ctx *builderContext) error {
	return imageContext(ctx, func(ctx *builderContext) error {
		runner := "camel-k-integration-" + defaults.Version + "-runner"
		ctx.Artifacts = []v1.Artifact{
			{
				ID:       runner,
				Location: QuarkusRuntimeSupport(ctx.Catalog.GetCamelQuarkusVersion()).TargetDirectory(ctx.Path, runner),
				Target:   runner,
			},
		}
		ctx.SelectedArtifacts = ctx.Artifacts
		return nil
	})
}

func executableDockerfile(ctx *builderContext) error {
	// #nosec G202
	dockerfile := []byte(`
		FROM ` + ctx.BaseImage + `
		WORKDIR ` + DeploymentDir + `
		COPY --chown=nonroot:root . ` + DeploymentDir + `
		USER nonroot
	`)

	err := os.WriteFile(filepath.Join(ctx.Path, ContextDir, "Dockerfile"), dockerfile, io.FilePerm400)
	if err != nil {
		return err
	}

	return nil
}

func standardImageContext(ctx *builderContext) error {
	return imageContext(ctx, func(ctx *builderContext) error {
		ctx.SelectedArtifacts = ctx.Artifacts

		return nil
	})
}

func jvmDockerfile(ctx *builderContext) error {
	// #nosec G202
	dockerfile := []byte(`
		FROM ` + ctx.BaseImage + `
		ADD . ` + DeploymentDir + `
		USER 1000
	`)

	err := os.WriteFile(filepath.Join(ctx.Path, ContextDir, "Dockerfile"), dockerfile, io.FilePerm400)
	if err != nil {
		return err
	}

	return nil
}

func incrementalImageContext(ctx *builderContext) error {
	images, err := listPublishedImages(ctx)
	if err != nil {
		return err
	}

	return imageContext(ctx, func(ctx *builderContext) error {
		ctx.SelectedArtifacts = ctx.Artifacts

		bestImage, commonLibs := findBestImage(images, ctx.Artifacts)
		if bestImage.Image != "" {

			log.Infof("Selected %s as base image for %s", bestImage.Image, ctx.Build.Name)
			ctx.BaseImage = bestImage.Image
			ctx.SelectedArtifacts = make([]v1.Artifact, 0)

			for _, entry := range ctx.Artifacts {
				if _, isCommon := commonLibs[entry.ID]; !isCommon {
					ctx.SelectedArtifacts = append(ctx.SelectedArtifacts, entry)
				}
			}
		}

		return nil
	})
}

func imageContext(ctx *builderContext, selector artifactsSelector) error {
	err := selector(ctx)
	if err != nil {
		return err
	}

	contextDir := filepath.Join(ctx.Path, ContextDir)

	err = os.MkdirAll(contextDir, io.FilePerm755)
	if err != nil {
		return err
	}

	for _, entry := range ctx.SelectedArtifacts {
		_, err := util.CopyFile(entry.Location, filepath.Join(contextDir, entry.Target))
		if err != nil {
			return err
		}
	}

	for _, entry := range ctx.Resources {
		filePath, fileName := path.Split(entry.Target)
		fullPath := filepath.Join(contextDir, filePath, fileName)
		if err := util.WriteFileWithContent(fullPath, entry.Content); err != nil {
			return err
		}
	}

	return nil
}

func listPublishedImages(context *builderContext) ([]v1.IntegrationKitStatus, error) {
	excludeNativeImages, err := labels.NewRequirement(v1.IntegrationKitLayoutLabel, selection.NotEquals, []string{
		v1.IntegrationKitLayoutNativeSources,
	})
	if err != nil {
		return nil, err
	}

	options := []ctrl.ListOption{
		ctrl.InNamespace(context.Namespace),
		ctrl.MatchingLabels{
			v1.IntegrationKitTypeLabel:           v1.IntegrationKitTypePlatform,
			kubernetes.CamelLabelRuntimeVersion:  context.Catalog.Runtime.Version,
			kubernetes.CamelLabelRuntimeProvider: string(context.Catalog.Runtime.Provider),
		},
		ctrl.MatchingLabelsSelector{
			Selector: labels.NewSelector().Add(*excludeNativeImages),
		},
	}

	list := v1.NewIntegrationKitList()
	err = context.Client.List(context.C, &list, options...)
	if err != nil {
		return nil, err
	}

	images := make([]v1.IntegrationKitStatus, 0)
	for _, kit := range list.Items {
		// Discard non ready kits
		if kit.Status.Phase != v1.IntegrationKitPhaseReady {
			continue
		}
		// Discard kits with a different root hierarchy
		// context.BaseImage should still contain the root base image at this stage
		if kit.Status.RootImage != context.BaseImage {
			continue
		}
		images = append(images, kit.Status)
	}
	return images, nil
}

func findBestImage(images []v1.IntegrationKitStatus, artifacts []v1.Artifact) (v1.IntegrationKitStatus, map[string]bool) {
	var bestImage v1.IntegrationKitStatus

	if len(images) == 0 {
		return bestImage, nil
	}

	requiredLibs := make(map[string]string, len(artifacts))
	for _, entry := range artifacts {
		requiredLibs[entry.ID] = entry.Checksum
	}

	bestImageCommonLibs := make(map[string]bool)

	for _, image := range images {
		nonLibArtifacts := 0
		common := make(map[string]bool)

		for _, artifact := range image.Artifacts {
			// the application artifacts should not be considered as dependencies for image reuse
			// otherwise, checksums would never match and we would always use the root image
			if !strings.HasPrefix(artifact.Target, "dependencies/lib") {
				nonLibArtifacts++
				continue
			}

			// If the Artifact's checksum is not defined we can't reliably determine if for some
			// reason the artifact has been changed but not the ID (as example for snapshots or
			// other generated jar) thus we do not take this artifact into account.
			if artifact.Checksum == "" {
				continue
			}
			if requiredLibs[artifact.ID] == artifact.Checksum {
				common[artifact.ID] = true
			}
		}

		numCommonLibs := len(common)
		surplus := len(image.Artifacts) - numCommonLibs - nonLibArtifacts

		if surplus > 0 {
			// the base image cannot have extra libs that we don't need
			continue
		}

		if numCommonLibs >= len(bestImageCommonLibs) {
			bestImage = image
			bestImageCommonLibs = common
		}
	}

	return bestImage, bestImageCommonLibs
}
