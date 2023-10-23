//go:build mocks

package mocks

//go:generate mockery --output . --filename ./project_repository_mock.go 	--dir .. --name Project
//go:generate mockery --output . --filename ./scenario_repository_mock.go 	--dir .. --name Scenario
//go:generate mockery --output . --filename ./run_repository_mock.go 		--dir .. --name Run
