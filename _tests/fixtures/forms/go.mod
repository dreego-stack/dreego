module forms

go 1.27

require (
	github.com/dreego-stack/dreego/adapter/ssr v0.8.0
	github.com/dreego-stack/dreego/core v0.8.0
)

replace github.com/dreego-stack/dreego => ../../..

replace github.com/dreego-stack/dreego/core => ../../../core

replace github.com/dreego-stack/dreego/adapter/ssr => ../../../adapter/ssr
