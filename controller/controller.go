package controller

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"net/http"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"fmt"
	"github.com/gaia-adm/mr-burns/common"
)

type Controller struct {
	taskIdToTask map[string]Task
	docker       DockerManager
	sleep        func()
	getResults   func(string) string
	publish      func(string) error
	stop         func() bool
}

func NewController(dockerManager DockerManager) Controller {

	return Controller{
		taskIdToTask: map[string]Task{},
		docker: dockerManager,
		sleep: func() {
			time.Sleep(time.Second * 10)
		},
		getResults: getTestResults,
		publish: publishResults,
		stop: func() bool {
			return false
		},
	}
}

func (controller Controller) Start() {

	controller.initialize()
	for !controller.stop() && !controller.isFinished() {
		controller.startWaitingTasks()
		controller.sleep()
	}
}

func (controller Controller) isFinished() bool {

	ret := true
	for _, currTask := range controller.taskIdToTask {
		if TASK_STATE_DONE != currTask.State {
			ret = false
			break
		}
	}

	return ret
}

func (controller Controller) startWaitingTasks() {

	for _, currTask := range controller.taskIdToTask {
		if TASK_STATE_WAITING == currTask.State &&
		currTask.NextRuntimeMillisecond < common.GetTimeNowMillisecond() {
			controller.startContainer(currTask)
		}
	}
}

func (controller Controller) initialize() {

	images, err := controller.docker.GetImages()
	if err != nil {
		log.Infof("Failed to get docker images. Error: %v", err)
	}

	for _, currImage := range images {
		controller.taskIdToTask[currImage.ID] = newTask(currImage)
	}
}

func (controller Controller) startContainer(task Task) {

	task.State = TASK_STATE_RUNNING
	controller.update(task)
	go func() {
		image := task.Data.(docker.APIImages)
		container := getContainerName(image)
		testResultsFilePath, err := controller.docker.RunTests(image, container)
		if err != nil {
			log.Infof("Error while trying to run tests from image: %s. Error: %v", image, err)
		}else {
			controller.publish(controller.getPublishData(testResultsFilePath, image, container))
		}
		controller.setTaskNextRunningTime(task)
	}()
}

func (controller Controller) update(task Task) {

	controller.taskIdToTask[task.ID] = task
}

func (controller Controller) setTaskNextRunningTime(task Task) {

	image := task.Data.(docker.APIImages)
	imageInterval := controller.docker.GetLabelImageRunningInterval(image)
	if (len(imageInterval) > 0) {
		interval, _ := strconv.ParseInt(imageInterval, 10, 64)
		task.NextRuntimeMillisecond = common.GetTimeNowMillisecond() + interval
		task.State = TASK_STATE_WAITING
	} else {
		task.NextRuntimeMillisecond = 0
		task.State = TASK_STATE_DONE
	}
	controller.update(task)
	log.Infof("Finish running image: %s (Tags: %s). Next run time: %d", image.ID, image.RepoTags, task.NextRuntimeMillisecond)
}

func (controller Controller) getPublishData(testResultsFilePath string, image docker.APIImages, container string) string {

	testResults := controller.getResults(testResultsFilePath)
	if len(testResults) == 0 {
		testResults, _ = controller.docker.GetContainerLogs(container)
	}
	testDesc := controller.docker.GetLabelImageDesc(image)
	if len(testDesc) > 0 {
		testResults = fmt.Sprintf("%s\n%s\n%s", testDesc, image.RepoTags[0], testResults)
	}
	log.Infof("Container: %s, test results: %s", container, testResults)

	return testResults
}

func getTestResults(testResultsFilePath string) string {

	testResults, err := ioutil.ReadFile(testResultsFilePath)
	if err != nil {
		log.Error("Failed to read test results file.", err, testResultsFilePath)
		testResults = []byte("")
	}

	return string(testResults)
}

func publishResults(data string) error {

	req, err := http.NewRequest("POST", "http://distributor-link:8000", bytes.NewBufferString(data))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Error("Failed to POST container test results", data, err)
		return err
	}
	defer response.Body.Close()

	return nil
}

func newTask(image docker.APIImages) Task {

	return Task{
		ID: image.ID,
		NextRuntimeMillisecond: -1,
		State: TASK_STATE_WAITING,
		Data: image,
	}
}

func getContainerName(image docker.APIImages) string {

	ret := image.RepoTags[0]
	ret = strings.Replace(ret, "/", ".", -1)

	return strings.Replace(ret, ":", ".", -1)
}