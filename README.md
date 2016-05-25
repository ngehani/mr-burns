# Mr. Burns - Micro services testing framework
Mr. Burns is a simple service for running Docker containers for integration tests.
## Before running Mr. Burns
Add to your integration test's _Dockerfile_ the follow labels:

* Label **_test_** (mandatory) tells Mr. Burns that this image is an image that Mr. Burns should run.

* Label **_test.container.settings_** (optional) is a JSON that tells Mr. Burns how to run your integration test Docker image.
For more details, see [Docker Remote API](https://docs.docker.com/reference/api/docker_remote_api_v1.16/#create-a-container).

* Label **_test.results.dir_** (mandatory) is where your tests results are stored.
When you finish running your tests, Mr. Burns will copy this folder into the host's folder under the directory /tmp/test-results/<container-name_timestamp>.

* Label **_test.results.file_** (mandatory) is the test results file to publish for applications such as HPE ALM, HPE NGA, and Slack.

* Label **_test.run.interval_** (optional) if you want to run your container recurrently, set the time interval to be between when it finishes running to when it starts over.

* Label **_test.run.events_** (optional) if you want to run your container when an Docker event occurs, you should define the events on which it should run.
For example: You want to run [docker-bench-test](https://github.com/gaia-adm/docker-bench-test/) (tests that contain common best practices for deploying Docker containers in production), each time a host runs or creates a new Docker container. In this case, set _test.run.events=create,run_.

* Label **_test.use.latest_** (optional, default false) if you set this option to true, Mr. Burns will use the latest image (Docker will pull it if necessary).

#####Example:
```
LABEL test=true
LABEL test.container.settings={\"Config\":{\"Env\":[\"gaiaUrl=master.gaiahub.io\"]}}\"
LABEL test.results.dir=/src/results
LABEL test.results.file=TestSuite.txt
LABEL test.run.interval=300000
```
## Running Mr. Burns
By default, Mr. Burns looks for all hosts' images that contains _test=true_ label in the _Dockerfile_ and runs them as a Docker container.

If you want to run only specific images, add `TEST_IMAGES=[image-name-1,image-name2...]` as an environment variable. 
### Running Mr. Burns inside a Docker container
If you use Mr. Burns Docker [image](https://hub.docker.com/r/gaiaadm/mr-burns/) on Linux, use the following command:
```bash
docker run -d --name burns --log-driver=json-file -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp gaiaadm/mr-burns
```
### Running Mr. Burns on CoreOS cluster
If you are running a CoreOS cluster, you can use the `fleetctl` command to deploy Mr. Burns service file on every CoreOS cluster node.
```
$ fleetctl start mr-burns.service
```
