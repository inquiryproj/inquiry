//go:build mocks

package mocks

//go:generate mockery --output . --filename echo_mock.go --srcpkg=github.com/labstack/echo/v4 --name=Context

//go:generate mockery --output . --filename ./api_key_repository_mock.go 	--dir .. --name APIKeyRepository

//go:generate mockery --output . --filename ./api_server_mock.go 	--dir ../api --name ServerInterface

//go:generate mockery --output . --filename ./runnable_mock.go 	--dir .. --name Runnable
