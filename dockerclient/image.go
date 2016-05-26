package dockerclient

import "github.com/fsouza/go-dockerclient"

const (
	LABEL_TEST = "test"
	LABEL_TEST_RUN_INTERVAL = "test.run.interval"
	LABEL_TEST_RESULTS_DIR = "test.results.dir"
	LABEL_TEST_RESULTS_FILE = "test.results.file"
	LABEL_TEST_CONTAINER_SETTINGS = "test.container.settings"
	LABEL_TEST_ENV_FILE = "test.environment.file"
)

func RunInterval(image docker.APIImages) string {

	return getLabel(image, LABEL_TEST_RUN_INTERVAL)
}

func ResultsDir(image docker.APIImages) string {

	return getLabel(image, LABEL_TEST_RESULTS_DIR)
}

func ResultsFile(image docker.APIImages) string {

	return getLabel(image, LABEL_TEST_RESULTS_FILE)
}

func ContainerSettings(image docker.APIImages) string {

	return getLabel(image, LABEL_TEST_CONTAINER_SETTINGS)
}

func EnvFile(image docker.APIImages) string {

	return getLabel(image, LABEL_TEST_ENV_FILE)
}

func getLabel(image docker.APIImages, label string) string {

	return image.Labels[label]
}