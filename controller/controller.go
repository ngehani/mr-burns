package controller

import (
	"time"
	"log"
	"github.com/fsouza/go-dockerclient"
	"net/http"
	"bytes"
	"io/ioutil"
	"strconv"
	"syscall"
	"strings"
)

type Controller struct {
	taskIdToTask map[string]Task
	docker       DockerManager
	sleep        func()
	publish      func(string) error
	stop         func() bool
}

func NewController(dockerManager DockerManager) Controller {

	return Controller{
		taskIdToTask: map[string]Task{},
		docker: dockerManager,
		sleep: controllerSleep,
		publish: publishResults,
		stop: controllerStop,
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
		currTask.NextRuntimeMillisecond < getTimeNowMillisecond() {
			controller.startContainer(&currTask)
		}
	}
}

func (controller Controller) initialize() {

	images, err := controller.docker.GetImages()
	if err != nil {
		log.Printf("Failed to get docker images. Error:%v", err)
	}
	for _, currImage := range images {
		controller.taskIdToTask[currImage.ID] = newTask(currImage)
	}
}

func (controller Controller) startContainer(task *Task) {

	task.State = TASK_STATE_RUNNING
	go func() {
		image := task.Data.(docker.APIImages)
		testResultsFilePath, _ := controller.docker.RunTests(image, getContainerName(image))
		controller.publish(testResultsFilePath)
		controller.setTaskNextRunningTime(task)
	}()
}

func (controller Controller) setTaskNextRunningTime(task *Task) {

	image := task.Data.(docker.APIImages)
	interval, _ := strconv.ParseInt(controller.docker.GetImageRunningInterval(image), 10, 64)
	if (interval > 0) {
		task.NextRuntimeMillisecond = getTimeNowMillisecond() + interval
		task.State = TASK_STATE_WAITING
	} else {
		task.NextRuntimeMillisecond = 0
		task.State = TASK_STATE_DONE
	}
	log.Printf("Finish running image: %s (Tags: %s). Next run time: %d", image.ID, image.RepoTags, task.NextRuntimeMillisecond)
}

func publishResults(testResultsFilePath string) error {

	testResults, err := ioutil.ReadFile(testResultsFilePath)
	if err != nil {
		log.Printf("Failed to read test results file. File: %s Error: %v", testResultsFilePath, err)
		return err
	}
	req, err := http.NewRequest("POST", "http://distributor-link:8000", bytes.NewBuffer(testResults))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Print("Failed to POST container test results", err)
		return err
	}
	defer response.Body.Close()

	return nil
}

func newTask(image docker.APIImages) Task {

	return Task{
		NextRuntimeMillisecond: -1,
		State: TASK_STATE_WAITING,
		Data: image,
	}
}

func getTimeNowMillisecond() int64 {

	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)

	return (int64(tv.Sec) * 1e3 + int64(tv.Usec) / 1e3)
}

func controllerSleep() {

	time.Sleep(time.Second * 10)
}

func controllerStop() bool {

	return false
}

func getContainerName(image docker.APIImages) string {

	ret := image.RepoTags[0]
	ret = strings.Replace(ret, "/", "_", -1)

	return strings.Replace(ret, ":", "_", -1)
}