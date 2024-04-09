package model

type Metadata struct {
	Name   string
	Lables []map[string]string
}

type PodSpec struct {
	Containers []Container
}

type Container struct {
	Name  string
	Image string
}

type Pod struct {
	Metadata
	PodSpec
}
