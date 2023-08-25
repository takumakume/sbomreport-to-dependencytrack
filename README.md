# sbomreport-to-dependencytrack

Tool to send trivy-operator's SBOM Report to Dependency Track,
which can receive webhooks from Stdin and TrivyOperator and send them to Dependency Track.
Project and tags can be generated with templates using SBOM Report values.

```shell
# from stdin
$ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack

# from webhook
$ sbomreport-to-dependencytrack server --port 80

# set project name, version and tags
#  - using go template with sprig functions 
#  - delimiter: "[[" "]]" (no conflict with helm template)
#  - ".sbomReport" variable: the root of the SBOM Report
$ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack \
    --base_url http://localhost:8081 \
    --api-key 1234567890 \
    --project-name "[[ .sbomReport.report.artifact.repository ]]"
    --project-version "[[ .sbomReport.report.artifact.tag ]]"
    --project-tags "tag1,kube_cluster_name:production,kube_namespace:[[ .sbomReport.report.metadaga.namespace ]]"

# set by environment variables
$ kubectl get sbom hoge -o json | \
    DT_BASE_URL=http://localhost:8081 \
    DT_API_KEY=1234567890 \
    DT_PROJECT_NAME="[[ .sbomReport.report.artifact.repository ]]" \
    DT_PROJECT_VERSION="[[ .sbomReport.report.artifact.tag ]]" \
    DT_PROJECT_TAGS="tag1,kube_cluster_name:production,kube_namespace:[[ .sbomReport.report.metadaga.namespace ]]" \
    sbomreport-to-dependencytrack
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
