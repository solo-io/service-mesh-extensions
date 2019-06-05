---
title: Contributing
menuTitle: Contributing
weight: 1
---

The Service Mesh Hub is a space designed for service mesh application developers to share their work and to allow 
for easy installation and management of said applications.

Contributions to the Service Mesh Hub come in the form of extension specs. More details on authoring these 
extensions can be found [here](../extensions).

An explanation on the basic structure of the hub repo can be found [here](https://github.com/solo-io/service-mesh-hub), 
this doc will explain in slightly more detail the structure of an individual extension folder.

#### Extensions

```
gloo
├── description.md
├── overlays
│   ├── supergloo
│   │   ├── kustomization.yaml
│   │   └── mesh-ingress.yaml
│   └── vanilla
│       ├── custom-resources.yaml
│       └── kustomization.yaml
├── spec.yaml
└── test
    ├── gloo_suite_test.go
    └── gloo_test.go
```

Above is the structure of the v1 Gloo extension. There are four main sections to an extension.
1) description.md
    * the long description of an extension, rendered on the extension's install page
1) overlays
	* kustomization style yaml which can be applied to an extension manifest after it has been rendered
2) spec.yaml
	* the main application spec for a given extension
3) test
	* unit tests to check the correctness of a given application spec.

For more detailed information on overlays and spec.yaml please see [here](../extensions).

##### Testing
For testing we use the go testing library `ginkgo`. The application spec tests are unit style tests, so they are
meant to test the correctness of the yaml, and the resources they create. They are NOT meant to test how they 
will function once applied/installed to a live kubernetes cluster. We have created a testing library to aid in 
this process which can be found [here](https://github.com/solo-io/go-utils/tree/master/manifesttestutils). 
For some examples of our expectations surrounding application spec tests see 
[`gloo`](https://github.com/solo-io/service-mesh-hub/blob/master/extensions/v1/gloo/test/gloo_test.go) and 
[`kiali`](https://github.com/solo-io/service-mesh-hub/blob/master/extensions/v1/kiali/test/kiali_test.go).

#### Pull Requests
Here at [solo.io](https://www.solo.io/) we envision the Service Mesh Hub will be a space for all 
application developers to easily share their service mesh application with the wider communnity, and thus
we are thrilled to accept community Pull Requests. In order to submit an application to the Service Mesh Hub simply
open a Pull Request in this repo containing your full, valid application spec directory, and a solo team member will
review it before it can be merged in.
