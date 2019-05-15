---
title: Authoring Extensions
menuTitle: Authoring Extensions
weight: 2
---

This document will guide you through the process of publishing your own mesh extension to the `Service Mesh Hub`.

--------
For any question you might have that is not addressed here, please join our [community Slack channel](https://slack.solo.io/).

--------

## Glossary
- `Extension directory`: the directory in `extensions/v1` that contains the specification for your extension.
- `Extension Spec`: the `.yaml` file at the base of your `Extension directory` that described your mesh extension.
- `Versioned Spec`: a section of the `Extension Spec` that represents a particular version of the extension.
- `Installation spec`: represents a location where the `Service Mesh Hub` can retrieve the information to install the basic version of an extension
- `Flavor`: each `Versioned Spec` can have one or more flavors. Flavors represent different variations of an extension. 
In most cases these will correspond to additional configuration specific to a particular mesh type or mesh version.

## Adding your extension
Extensions are added to the Service Mesh Hub by creating a folder in the `extensions/v1` directory. The name of this 
directory (from now on referred to as the `Extension directory`) must be equal to the name of your mesh extension. For a 
comprehensive explanation on how to add your extension to the `Service Mesh Hub`, check out our [contribution guide](../contributing).

The Extension Directory must always contain a file named `spec.yaml` at its root. This file, called the `Extension spec`, 
contains a complete description of your application. The rest of the content of your directory depends on the content of 
your `Extension spec`.

## The Extension Spec
Let's go over the most important attributes of the `Extension spec`. To help us along end we will define a sample mesh 
extension.

---
NOTE: The ultimate source of truth for the structure of the `Extension spec` is the 
[`.proto` file](https://github.com/solo-io/service-mesh-hub/blob/master/api/v1/registry.proto) located in the `/api` 
directory of our repository. Be sure to check it out!

---

This is the spec for the `strainer` extension:

```yaml
name: strainer
applicationCreator: solo.io
applicationMaintainer: solo.io
applicationProvider: solo.io
documentationUrl: https://docs.strainer.io/
logoUrl: https://some.domain/path/to/an/image
shortDescription: Strainer is a demo application to illustrate the contents of an Extension Spec.
longDescription: |
  This attribute should contain a comprehensive description of the extension.
repositoryUrl: github.com/solo-io/strainer
versions:
- version: "1.0.0"
  githubChart:
    org: solo.io
    repo: strainer
    ref: "v1.0.0"
    directory: installation/chart
  datePublished: "2019-05-15T00:00:00Z"
  valuesYaml: |-
    rbac:
      create: true
    mesh:
      name: {{ .MeshRef.Name }}
      namespace: {{ .MeshRef.Namespace }}
  flavors:
  - name: istio
    description: "Configure strainer to leverage the features of your Istio mesh"
    requirementSets:
    - meshRequirement:
      meshType: ISTIO
        versions:
          maxVersion: "1.0.7"
          minVersion: "1.0.0"
    parameters:
    - name: prometheus-url
      description: |-
        The URL of the prometheus service that strainer uses to scrape metrics
      required: true
      default: "http://prometheus.{{ .MeshRef.Namespace }}:9090"
  - name: linkerd
    description: "Configure strainer to leverage the features of your Linkerd2 mesh"
    customizationLayers:
    - kustomize:
        github:
          org: solo-io
          repo: strainer
          ref: "v1.0.0"
          directory: installation/kustomize
        overlayPath: linkerd
    requirementSets:
    - meshRequirement:
      meshType: LINKERD
        versions:
          minVersion: "v18.7.1"
    parameters:
    - name: prometheus-url
      description: |-
        The URL of the prometheus service that strainer uses to scrape metrics
      required: true
      default: "http://prometheus.{{ .MeshRef.Namespace }}:9090"
    - name: write-namespace
      description: >
        An additional parameter that is required to apply the Linkerd customization
      default: "{{ .InstallNamespace }}"
      required: true
```

Let's look at the different sections individually.

### General extension information
The first part of the `Extension Spec` contains general information about the extension. These include i.a. links to the 
GitHub repository and documentation sites, information about the developer, and descriptions of the extension itself.

```yaml
name: strainer
applicationCreator: solo.io
applicationMaintainer: solo.io
applicationProvider: solo.io
documentationUrl: https://docs.strainer.io/
logoUrl: https://some.domain/path/to/an/image
shortDescription: Strainer is a demo application to illustrate the contents of an Extension Spec.
longDescription: |
  This attribute should contain a comprehensive description of the extension.
repositoryUrl: github.com/solo-io/strainer
```

The Service Mesh Hub UI will use this information to populate the entry for your extension.

### Versioned information
The top level `versions` attribute consists of an array of `VersionedApplicationSpec`s. Each one of these objects 
represents a distinct version of the extension.

Following is a description of the attributes of a `VersionedApplicationSpec` (from now `Versioned spec`)

#### Installation spec
Every `VersionedApplicationSpec` must specify a location from which the `Service Mesh Hub` will retrieve the installation 
manifest for this version of the extension. This can be done via three different top-level attributes.

##### 1. githubChart
Represents a Helm chart stored in a GitHub repository. For example, in the case of the `strainer` example above:

```yaml
githubChart:
  org: solo.io
  repo: strainer
  ref: "v1.0.0"
  directory: installation/chart
```

With this configuration, the `Service Mesh Hub` expects to find a Helm chart in the `installation/chart` directory of the 
`v1.0.0` ref (a branch, tag, or commit SHA) of the git repository located at `https://github.com/solo-io/strainer`.

##### 2. helmArchive
Represents a Helm chart archive (`gzip`ped tarball) stored at the given URI:

```yaml
helmArchive:
  uri: "https://storage.googleapis.com/my-bucket/strainer-1.0.0.tgz"
```

##### 3. manifestsArchive
Represents a `gzip`ped tarball containing a plain kubernetes `yaml` manifest stored at the given URI:

```yaml
manifestsArchive:
  uri: "https://storage.googleapis.com/my-bucket/strainer-manifest-1.0.0.tgz"
```

##### Default values
If your installation spec is a Helm chart, you can include values to be used during the rendering of the chart via the 
`valuesYaml` attribute:

```yaml
valuesYaml: |-
  rbac:
    create: true
  mesh:
    name: {{ .MeshRef.Name }}
    namespace: {{ .MeshRef.Namespace }}
```

The values can contain [go template actions](https://golang.org/pkg/text/template/) (text delimited by `{{` and `}}`). 
these placeholders will be resolved by the `Service Mesh Hub` before being passes to the Helm rendering engine.

See the [injected values](#Injected-values) section of this guide for a full list of the values that that are 
available for injection and how they are selected.

#### Flavors
A mesh extension might work only on certain mesh types and only an a subset of the available versions of the given 
mesh type. Furthermore, the `Installation spec` might require some tweaking based on which mesh and version the 
extension is being installed to. While Helm provides a powerful templating engine that allows you to define resources 
based on the provided values, conditional logic inside template files can quickly grow beyond maintainability when you 
have to account for several combinations of parameter values.

Extension `flavors` are meant to address this problem by making it simple to define and reason about the mesh-specific 
concerns of your extension. Each version of your extension *must* define at least one `flavor`. At a minimum, a flavor 
must have a `name` and a `description`, indicating that the given version of the extension will work on any mesh without 
further customization, but that's a boring scenario. Let us go through the rest of the (more interesting) attributes 
that make up a `flavor`.

##### Requirement sets
The `requirementSets` attribute defines sets of conditions that need to be satisfied for a `flavor` to be installable and 
consists of an array of `RequirementSet`s. A flavor is considered installable if at least one of the `requirementSets` 
is satisfied.

Currently, a requirement set has a single `meshRequirement` attribute, which defines mesh specific constraints. We are 
planning on adding more requirement types, e.g. on applications that have to be deployed to your cluster for an 
extension to be installed. Here are the requirements sets for the `istio` flavor of `strainer` sample app:

```yaml
requirementSets:
- meshRequirement:
    meshType: ISTIO
    versions:
      maxVersion: "1.0.7"
      minVersion: "1.0.0"
```

This means that the flavor can be installed only if you have Istio deployed to your cluster and the Istio version is 
greater than or equal to `1.0.0` and less than or equal to `1.0.7`.

The `minVersion`/`maxVersion` attributes are optional. If one is missing, the correspondent lower/upper boundary is 
removed. The `version` attribute itself can be omitted, in which case the flavor applies to all versions of the given 
mesh. This will for example be the case the flavor requires AWS App Mesh to be configured: since the App Mesh control 
plane is not versioned, we can simple define the requirement set as:

```yaml
requirementSets:
- meshRequirement:
    meshType: AWS_APP_MESH
```

As previously mentioned, if more than one requirement set is defined, the flavor is installable if any one of them is 
satisfied. A flavor with the following configuration can be installed either on istio or linkerd:

```yaml
requirementSets:
- meshRequirement:
    meshType: ISTIO
- meshRequirement:
    meshType: LINKERD
```

##### Customization layers
The `customizationLayers` attribute lets you define modifications to your installation that are specific to the 
given flavor. It consists of an array of `layers`. A `layer` represents a set of customizations that are logically 
related. Currently we support only one `layer` per flavor, but we envision having multiple predefined layers which can 
be "superimposed" and injected with different [parameters](#Parameters) to provide a powerful and flexible configuration 
mechanism. 

Layers are currently implemented using [kustomize](https://github.com/kubernetes-sigs/kustomize). Let's look at the 
`linkerd` flavor of `strainer` sample app:

```yaml
customizationLayers:
  - kustomize:
      github:
        org: solo-io
        repo: strainer
        ref: "v1.0.0"
        directory: installation/kustomize
      overlayPath: linkerd
```

A `kustomize `customization layer consists of a location and an overlay path.

The location refers to the root of a `kustomize` directory structure. This can be either a directory inside a GitHub 
repository or a remote archive containing the directory structure. The above yaml snippet represents the first case. A 
remote `kustomize` layer has the following form:

```yaml
customizationLayers:
  - kustomize:
      tgzArchive:
        uri: "https://storage.googleapis.com/my-bucket/strainer-1.0.0.tgz"
      overlayPath: linkerd
```

The `overlayPath` specifies the path to an [overlay](https://github.com/kubernetes-sigs/kustomize/blob/master/docs/glossary.md#overlay) 
inside the directory structure.

During installation of the extension, the `Service Mesh Hub` will configure the rendered Helm manifest to be the 
[base](https://github.com/kubernetes-sigs/kustomize/blob/master/docs/glossary.md#base) of the `kustomization`.

##### Parameters
Each flavor can define a set of parameters that define values that will be supplied to the execution of the 
customization layers. The `Service Mesh Hub` UI will dynamically generate input fields for each parameter when the user 
selects a flavor and include all the information provided in the `parametes` section.

Here is an example from the `linkerd` flavor of `strainer` sample app:

```yaml
parameters:
- name: prometheus-url
  description: |-
    The URL of the prometheus service that strainer uses to scrape metrics
  required: true
  default: "http://prometheus.{{ .MeshRef.Namespace }}:9090"
- name: write-namespace
  description: >
    An additional parameter that is required to apply the Linkerd customization
  default: "{{ .InstallNamespace }}"
  required: true
```

NOTE: the UI currently does not support parameters, so `default` values will always be used instead.

### Injected values
Check out the `ValuesInputs` object for values that are available during rendering of template actions in `valuesYaml`s
and flavor parameters.
