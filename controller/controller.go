package controller

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/gaia-adm/mr-burns/common"
	"fmt"
	"github.com/gaia-adm/mr-burns/dockerclient"
)

type Controller struct {
	taskIdToTask map[string]Task
	docker       DockerManager
	sleep        func()
	getResults   func(string) string
	publish      func(string) error
	stop         func() bool
}

func NewController(dockerManager DockerManager, publisher Publisher) Controller {

	return Controller{
		taskIdToTask: map[string]Task{},
		docker: dockerManager,
		sleep: func() {
			time.Sleep(time.Second * 10)
		},
		getResults: getTestResults,
		publish: publisher.Publish,
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
			log.Infof("Error while trying to run tests from image: %v. Error: %v", image, err)
		} else {
			controller.publish(controller.getPublishData(testResultsFilePath, image, container))
		}
		controller.setTaskNextRunningTime(task)
		log.Infof("Finish running container: %s, next run time: %v",
			container, common.MillisecondToTime(task.NextRuntimeMillisecond))
	}()
}

func (controller Controller) update(task Task) {

	controller.taskIdToTask[task.ID] = task
}

func (controller Controller) setTaskNextRunningTime(task Task) {

	image := task.Data.(docker.APIImages)
	imageInterval := dockerclient.RunInterval(image)
	if (len(imageInterval) > 0) {
		interval, _ := strconv.ParseInt(imageInterval, 10, 64)
		task.NextRuntimeMillisecond = common.GetTimeNowMillisecond() + interval
		task.State = TASK_STATE_WAITING
	} else {
		task.NextRuntimeMillisecond = 0
		task.State = TASK_STATE_DONE
	}
	controller.update(task)
}

func (controller Controller) getPublishData(testResultsFilePath string, image docker.APIImages, container string) string {

	var ret string
	testResults := controller.getResults(testResultsFilePath)
	if len(testResults) > 0 {
		ret = fmt.Sprintf("Container: %s\n%+v\n%s", container, image.RepoTags, testResults)
		log.Info(ret)
	} else {
		containerLogs, err := controller.docker.GetContainerLogs(container)
		if err != nil {
			ret = fmt.Sprintf("Empty container test results. Trying to fetch container's logs but failed, Error: %v (Container: %s %+v)", err, container, image.RepoTags)
			log.Error(ret)
		} else {
			ret = fmt.Sprintf("Container: %s\n%+v\nEmpty Test Results, Logs:\n%s", container, image.RepoTags, containerLogs)
			log.Info(ret)
		}
	}

	return ret
}

func getTestResults(testResultsFilePath string) string {

	testResults, err := ioutil.ReadFile(testResultsFilePath)
	if err != nil {
		log.Error("Failed to read test results file. ", err, testResultsFilePath)
		testResults = []byte("")
	}

	return string(testResults)
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