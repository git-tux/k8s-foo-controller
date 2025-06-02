* https://book.kubebuilder.io/cronjob-tutorial/api-design
* 
# Simple Kubernetes operator in golang

This repository demonstrates an simple operator which implements CRUD actions for a Custom Resource named **FooLabel**. For the development of the operator I will use the kubebuilder tool which creates the structure with the required files for building an operator using Golang.

## Initialize the project

```
kubebuilder init --domain foo.controller --repo github.com/git-tux/k8s-foo-controller
```
## Create an API

```
kubebuilder  create api --group foogroup --version v1 --kind FooLabel
```

## Create FooLabel struct

To create the required fields for the FooLabel kubernetes resource, I had to edit the file [foolabel_types.go](api/v1/foolabel_types.go). The FooLabel struct will be consisted only by one aattribute ( **label**) which contains the label that has to be added in the pods.

```go

type FooLabelSpec struct {
  // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
  // Important: Run "make" to regenerate code after modifying this file

  // Foo is an example field of FooLabel. Edit foolabel_types.go to remove/update
  Label string `json:"foo,omitempty"`
}

```

## Create controller logic

To add functionally in the controller, I edited the file [foolabel_controller.go](internal/controller/foolabel_controller.go).

## Create CRD

```
make install
```

## Create a foolabel instance

```
kubectl apply  -f config/samples/foogroup_v1_foolabel.yaml
```

## Build the operator docker image

```
echo $PAT | docker login ghcr.io -u git-tux --password-stdin
make docker-build docker-push IMG=ghcr.io/git-tux/foo-operator
```

## Deploy operator on kubernetes

```
make deploy IMG=ghcr.io/git-tux/foo-operator

```
