package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/neuralmux/devpod-provider-upcloud/features/step_definitions"
)

var opts = godog.Options{
	Format:   "pretty",
	Paths:    []string{"features"},
	Output:   colors.Colored(os.Stdout),
	TestingT: &testing.T{},
}

func TestFeatures(t *testing.T) {
	opts.TestingT = t

	status := godog.TestSuite{
		Name:                 "UpCloud Provider",
		ScenarioInitializer: step_definitions.InitializeScenario,
		TestSuiteInitializer: step_definitions.InitializeTestSuite,
		Options:              &opts,
	}.Run()

	if status != 0 {
		t.Fail()
	}
}