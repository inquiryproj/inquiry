//go:build mocks

package mocks

//go:generate mockery --output . --filename ./project_service_mock.go 	--dir .. --name Project
//go:generate mockery --output . --filename ./scenario_service_mock.go 	--dir .. --name Scenario
//go:generate mockery --output . --filename ./runner_service_mock.go 	--dir .. --name Runner
