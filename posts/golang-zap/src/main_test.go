package main

import "testing"

func Test_dev(t *testing.T) {
	dev()
}

func Test_test(t *testing.T) {
	test()
}

func Test_prod(t *testing.T) {
	prod()
}

func Test_devWithConfig(t *testing.T) {
	devWithConfig()
}

func Test_devWithConfig2(t *testing.T) {
	devWithConfigTimeFormat()
}

func Test_devWithField(t *testing.T) {
	devWithField()
}

func Test_devWithColor(t *testing.T) {
	devWithColor()
}

func Test_devWithCustomEncoder(t *testing.T) {
	devWithCustomEncoder()
}

func Test_devWithGlobal(t *testing.T) {
	devWithGlobal()
}
