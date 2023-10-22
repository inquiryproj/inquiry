// Code generated by mockery v2.19.0. DO NOT EDIT.

package mocks

import (
	io "io"
	multipart "mime/multipart"
	http "net/http"
	url "net/url"

	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// Context is an autogenerated mock type for the Context type
type Context struct {
	mock.Mock
}

// Attachment provides a mock function with given fields: file, name
func (_m *Context) Attachment(file string, name string) error {
	ret := _m.Called(file, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(file, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Bind provides a mock function with given fields: i
func (_m *Context) Bind(i interface{}) error {
	ret := _m.Called(i)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Blob provides a mock function with given fields: code, contentType, b
func (_m *Context) Blob(code int, contentType string, b []byte) error {
	ret := _m.Called(code, contentType, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string, []byte) error); ok {
		r0 = rf(code, contentType, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Cookie provides a mock function with given fields: name
func (_m *Context) Cookie(name string) (*http.Cookie, error) {
	ret := _m.Called(name)

	var r0 *http.Cookie
	if rf, ok := ret.Get(0).(func(string) *http.Cookie); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Cookie)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Cookies provides a mock function with given fields:
func (_m *Context) Cookies() []*http.Cookie {
	ret := _m.Called()

	var r0 []*http.Cookie
	if rf, ok := ret.Get(0).(func() []*http.Cookie); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*http.Cookie)
		}
	}

	return r0
}

// Echo provides a mock function with given fields:
func (_m *Context) Echo() *echo.Echo {
	ret := _m.Called()

	var r0 *echo.Echo
	if rf, ok := ret.Get(0).(func() *echo.Echo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*echo.Echo)
		}
	}

	return r0
}

// Error provides a mock function with given fields: err
func (_m *Context) Error(err error) {
	_m.Called(err)
}

// File provides a mock function with given fields: file
func (_m *Context) File(file string) error {
	ret := _m.Called(file)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(file)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FormFile provides a mock function with given fields: name
func (_m *Context) FormFile(name string) (*multipart.FileHeader, error) {
	ret := _m.Called(name)

	var r0 *multipart.FileHeader
	if rf, ok := ret.Get(0).(func(string) *multipart.FileHeader); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*multipart.FileHeader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FormParams provides a mock function with given fields:
func (_m *Context) FormParams() (url.Values, error) {
	ret := _m.Called()

	var r0 url.Values
	if rf, ok := ret.Get(0).(func() url.Values); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(url.Values)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FormValue provides a mock function with given fields: name
func (_m *Context) FormValue(name string) string {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *Context) Get(key string) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// HTML provides a mock function with given fields: code, html
func (_m *Context) HTML(code int, html string) error {
	ret := _m.Called(code, html)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(code, html)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HTMLBlob provides a mock function with given fields: code, b
func (_m *Context) HTMLBlob(code int, b []byte) error {
	ret := _m.Called(code, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []byte) error); ok {
		r0 = rf(code, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Handler provides a mock function with given fields:
func (_m *Context) Handler() echo.HandlerFunc {
	ret := _m.Called()

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// Inline provides a mock function with given fields: file, name
func (_m *Context) Inline(file string, name string) error {
	ret := _m.Called(file, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(file, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsTLS provides a mock function with given fields:
func (_m *Context) IsTLS() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsWebSocket provides a mock function with given fields:
func (_m *Context) IsWebSocket() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// JSON provides a mock function with given fields: code, i
func (_m *Context) JSON(code int, i interface{}) error {
	ret := _m.Called(code, i)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, interface{}) error); ok {
		r0 = rf(code, i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// JSONBlob provides a mock function with given fields: code, b
func (_m *Context) JSONBlob(code int, b []byte) error {
	ret := _m.Called(code, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []byte) error); ok {
		r0 = rf(code, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// JSONP provides a mock function with given fields: code, callback, i
func (_m *Context) JSONP(code int, callback string, i interface{}) error {
	ret := _m.Called(code, callback, i)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string, interface{}) error); ok {
		r0 = rf(code, callback, i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// JSONPBlob provides a mock function with given fields: code, callback, b
func (_m *Context) JSONPBlob(code int, callback string, b []byte) error {
	ret := _m.Called(code, callback, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string, []byte) error); ok {
		r0 = rf(code, callback, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// JSONPretty provides a mock function with given fields: code, i, indent
func (_m *Context) JSONPretty(code int, i interface{}, indent string) error {
	ret := _m.Called(code, i, indent)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, interface{}, string) error); ok {
		r0 = rf(code, i, indent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Logger provides a mock function with given fields:
func (_m *Context) Logger() echo.Logger {
	ret := _m.Called()

	var r0 echo.Logger
	if rf, ok := ret.Get(0).(func() echo.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.Logger)
		}
	}

	return r0
}

// MultipartForm provides a mock function with given fields:
func (_m *Context) MultipartForm() (*multipart.Form, error) {
	ret := _m.Called()

	var r0 *multipart.Form
	if rf, ok := ret.Get(0).(func() *multipart.Form); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*multipart.Form)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NoContent provides a mock function with given fields: code
func (_m *Context) NoContent(code int) error {
	ret := _m.Called(code)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Param provides a mock function with given fields: name
func (_m *Context) Param(name string) string {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ParamNames provides a mock function with given fields:
func (_m *Context) ParamNames() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// ParamValues provides a mock function with given fields:
func (_m *Context) ParamValues() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Path provides a mock function with given fields:
func (_m *Context) Path() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// QueryParam provides a mock function with given fields: name
func (_m *Context) QueryParam(name string) string {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// QueryParams provides a mock function with given fields:
func (_m *Context) QueryParams() url.Values {
	ret := _m.Called()

	var r0 url.Values
	if rf, ok := ret.Get(0).(func() url.Values); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(url.Values)
		}
	}

	return r0
}

// QueryString provides a mock function with given fields:
func (_m *Context) QueryString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RealIP provides a mock function with given fields:
func (_m *Context) RealIP() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Redirect provides a mock function with given fields: code, _a1
func (_m *Context) Redirect(code int, _a1 string) error {
	ret := _m.Called(code, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(code, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Render provides a mock function with given fields: code, name, data
func (_m *Context) Render(code int, name string, data interface{}) error {
	ret := _m.Called(code, name, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string, interface{}) error); ok {
		r0 = rf(code, name, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Request provides a mock function with given fields:
func (_m *Context) Request() *http.Request {
	ret := _m.Called()

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func() *http.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	return r0
}

// Reset provides a mock function with given fields: r, w
func (_m *Context) Reset(r *http.Request, w http.ResponseWriter) {
	_m.Called(r, w)
}

// Response provides a mock function with given fields:
func (_m *Context) Response() *echo.Response {
	ret := _m.Called()

	var r0 *echo.Response
	if rf, ok := ret.Get(0).(func() *echo.Response); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*echo.Response)
		}
	}

	return r0
}

// Scheme provides a mock function with given fields:
func (_m *Context) Scheme() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Set provides a mock function with given fields: key, val
func (_m *Context) Set(key string, val interface{}) {
	_m.Called(key, val)
}

// SetCookie provides a mock function with given fields: cookie
func (_m *Context) SetCookie(cookie *http.Cookie) {
	_m.Called(cookie)
}

// SetHandler provides a mock function with given fields: h
func (_m *Context) SetHandler(h echo.HandlerFunc) {
	_m.Called(h)
}

// SetLogger provides a mock function with given fields: l
func (_m *Context) SetLogger(l echo.Logger) {
	_m.Called(l)
}

// SetParamNames provides a mock function with given fields: names
func (_m *Context) SetParamNames(names ...string) {
	_va := make([]interface{}, len(names))
	for _i := range names {
		_va[_i] = names[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// SetParamValues provides a mock function with given fields: values
func (_m *Context) SetParamValues(values ...string) {
	_va := make([]interface{}, len(values))
	for _i := range values {
		_va[_i] = values[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// SetPath provides a mock function with given fields: p
func (_m *Context) SetPath(p string) {
	_m.Called(p)
}

// SetRequest provides a mock function with given fields: r
func (_m *Context) SetRequest(r *http.Request) {
	_m.Called(r)
}

// SetResponse provides a mock function with given fields: r
func (_m *Context) SetResponse(r *echo.Response) {
	_m.Called(r)
}

// Stream provides a mock function with given fields: code, contentType, r
func (_m *Context) Stream(code int, contentType string, r io.Reader) error {
	ret := _m.Called(code, contentType, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string, io.Reader) error); ok {
		r0 = rf(code, contentType, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// String provides a mock function with given fields: code, s
func (_m *Context) String(code int, s string) error {
	ret := _m.Called(code, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(code, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Validate provides a mock function with given fields: i
func (_m *Context) Validate(i interface{}) error {
	ret := _m.Called(i)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// XML provides a mock function with given fields: code, i
func (_m *Context) XML(code int, i interface{}) error {
	ret := _m.Called(code, i)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, interface{}) error); ok {
		r0 = rf(code, i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// XMLBlob provides a mock function with given fields: code, b
func (_m *Context) XMLBlob(code int, b []byte) error {
	ret := _m.Called(code, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []byte) error); ok {
		r0 = rf(code, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// XMLPretty provides a mock function with given fields: code, i, indent
func (_m *Context) XMLPretty(code int, i interface{}, indent string) error {
	ret := _m.Called(code, i, indent)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, interface{}, string) error); ok {
		r0 = rf(code, i, indent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewContext interface {
	mock.TestingT
	Cleanup(func())
}

// NewContext creates a new instance of Context. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewContext(t mockConstructorTestingTNewContext) *Context {
	mock := &Context{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
