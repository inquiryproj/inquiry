//go:build mocks

package mocks

//go:generate mockery --output . --filename echo_mock.go --srcpkg=github.com/labstack/echo/v4 --name=Context
