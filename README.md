# sbomreport-to-dependencytrack

```shell
Send Aqua Security Trivy Operator's SBOM Report to OWASP Dependency-Track.

two ways to use:

1. command line tool to receive JSON of SBOM Report from stdin

  $ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack

2. http server that receives JSON of SBOM Report from Trivy Operator webhook

  $ sbomreport-to-dependencytrack server --port 8080

Templates with the SBOM Report as a variable can be used for the following items to be registered in the Dependency-Track.

  * Project Name
  * Project Version
  * Project Tags
  
  $ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack \
      --base-url http://127.0.0.1:8081/ \
      --api-key 1234567890 \
      --project-name "[[.sbomReport.report.artifact.repository]]" \                  # e.g. "library/alpine"
      --project-version "[[.sbomReport.report.artifact.tag]]" \                      # e.g. "3.13.5"
      --project-tags "tag1,kube_namespace:[[.sbomReport.metadata.namespace]]" # e.g. ["tag1", "kube_namespace:default"]

  For template, go-template and sprig functions can be used.
  The delimiter of template is "[[" "]]". This is to avoid conflicts with other tools such as Helm.

Environment variables can be used instead of command line arguments, which may be useful when running on Kubernetes.

  $ kubectl get sbom hoge -o json | \
      DT_BASE_URL=http://127.0.0.1:8081 \
      DT_API_KEY=1234567890 \
      DT_PROJECT_NAME="[[.sbomReport.report.artifact.repository]]" \
      DT_PROJECT_VERSION="[[.sbomReport.report.artifact.tag]]" \
      DT_PROJECT_TAGS="tag1,kube_namespace:[[.sbomReport.metadata.namespace]]" \
      sbomreport-to-dependencytrack

Dependency-Track APK key permissions required:

  * BOM_UPLOAD 
  * PORTFOLIO_MANAGEMENT 
  * PROJECT_CREATION_UPLOAD 
  * VIEW_PORTFOLIO 
```

# Quick start

## Create API key for Dependency-Track

Admin > Access Management > Teams

Permissions:

  * BOM_UPLOAD 
  * PORTFOLIO_MANAGEMENT 
  * PROJECT_CREATION_UPLOAD 
  * VIEW_PORTFOLIO

## case1: command line tool to receive JSON of SBOM Report from stdin

install command

```shell
$ go install github.com/takumakume/sbomreport-to-dependencytrack@main
```

run command

```shell
$ cat testdata/v1alpha1.json | sbomreport-to-dependencytrack \
      --base-url http://<Dependency-Track API IP address>:<Port>/ \
      --api-key ********************************* \
      --project-name "[[.sbomReport.report.artifact.repository]]" \
      --project-version "[[.sbomReport.report.artifact.tag]]" \
      --project-tags "tag1,kube_namespace:[[.sbomReport.metadata.namespace]]"

2023/08/25 21:56:54 Uploading BOM: project library/alpine:latest
2023/08/25 21:56:54 Polling completion of upload BOM: project library/alpine:latest token aa5475a1-ff24-4402-b07b-c622733ea7ba
2023/08/25 21:56:55 BOM upload completed: project library/alpine:latest token aa5475a1-ff24-4402-b07b-c622733ea7ba
2023/08/25 21:56:55 Adding tags to project. project library/alpine:latest tags [tag1 kube_namespace:default]
```

## case2: http server that receives JSON of SBOM Report from Trivy Operator webhook

run server

```shell
$ docker run -p 8080:8080 \
  -e DT_BASE_URL=http://<Dependency-Track API IP address>:<Port>/ \
  -e DT_API_KEY=********************************* \
  -e DT_PROJECT_NAME="[[.sbomReport.report.artifact.repository]]" \
  -e DT_PROJECT_VERSION="[[.sbomReport.report.artifact.tag]]" \
  -e DT_PROJECT_TAGS="tag1,kube_namespace:[[.sbomReport.metadata.namespace]]" \
  -it docker.io/takumakume/sbomreport-to-dependencytrack:latest server

2023/08/25 13:05:41 Listening on :8080

# Run `curl localhost -X POST -d @testdata/v1alpha1.json`

2023/08/25 22:05:44 Uploading BOM: project library/alpine:latest
2023/08/25 22:05:44 Polling completion of upload BOM: project library/alpine:latest token 811585ae-39c9-402e-9e79-82e33a3d401d
2023/08/25 22:05:45 BOM upload completed: project library/alpine:latest token 811585ae-39c9-402e-9e79-82e33a3d401d
2023/08/25 22:05:45 Adding tags to project. project library/alpine:latest tags [tag1,kube_namespace:default]
```

# install

## kubernetes

### helm

```shell
helm repo add sbomreport-to-dependencytrack https://takumakume.github.io/sbomreport-to-dependencytrack/charts
helm repo update
helm install sbomreport-to-dependencytrack/sbomreport-to-dependencytrack

# render manifests
helm template sbomreport-to-dependencytrack/sbomreport-to-dependencytrack
```

main settings in values.yaml

```yaml
config:
  # Dependency Track API key secret name
  apiKeySecretName: sbomreport-to-dependencytrack
  
  # Dependency Track base URL
  baseUrl: "http://localhost:8081"

  # Dependency Track project name template
  projectName: "[[.sbomReport.report.artifact.repository]]"
  
  # Dependency Track project version template
  projectVersion: "[[.sbomReport.report.artifact.tag]]"
  
  # Dependency Track project tag template (comma separated)
  projectTags: ""
```

register webhook with trivy-operator

ref: https://aquasecurity.github.io/trivy-operator/latest/tutorials/integrations/webhook/

e.g.

```yaml
env:
  - name: OPERATOR_WEBHOOK_BROADCAST_URL
    value: http://sbomreport-to-dependencytrack.default.svc.cluster.local:8080/
```
