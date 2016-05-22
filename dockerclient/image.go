package dockerclient

import "github.com/fsouza/go-dockerclient"

const (
	LABEL_TEST = "test"
	LABEL_TEST_RUN_INTERVAL = "test.run.interval"
	LABEL_TEST_RESULTS_DIR = "test.results.dir"
	LABEL_TEST_RESULTS_FILE = "test.results.file"
	LABEL_TEST_CONTAINER_SETTINGS = "test.container.settings"
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

func getLabel(image docker.APIImages, label string) string {

	return image.Labels[label]
}