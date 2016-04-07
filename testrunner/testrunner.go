package testrunner

import (
	"github.com/gaia-adm/mr-burns/container"
)

func testContainersFilter(c container.Container) bool { return c.IsTest() }

func RunTestContainers(client container.burnsClient) error {
	//Get images with test label and create containers for them

	return nil
}


