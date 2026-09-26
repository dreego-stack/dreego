module importcheck

go 1.27

require (
	github.com/dreego-stack/dreego/adapter/ssr v0.10.9
	github.com/dreego-stack/dreego/core v0.10.9
	golang.org/x/text v0.22.0
)

replace github.com/dreego-stack/dreego => ../../..

replace github.com/dreego-stack/dreego/core => ../../../core

replace github.com/dreego-stack/dreego/adapter/ssr => ../../../adapter/ssr
