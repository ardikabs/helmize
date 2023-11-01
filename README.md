# Helmize is a KRM function for Helm on Kustomize with Glob support

The motivation for this project is quite simple. While [Helm integration](https://kubectl.docs.kubernetes.io/references/kustomize/builtins/#_helmchartinflationgenerator_) in Kustomize is already available, it falls short in scenarios where dynamic use of Glob is required to fetch all necessary values.yaml files. In such cases, it lacks support.
With the likelihood of [glob support](https://github.com/kubernetes-sigs/kustomize/issues/119) being added to Kustomize remaining uncertain for the foreseeable future, this project was initiated to address this limitation.

## Usage

The Release specification:

```bash
cat <<EOF > release.yaml
# release.yaml
apiVersion: toolkit.ardikabs.com/v1alpha1
kind: Release
metadata:
  name: service-a
  namespace: default
  annotations:
    config.kubernetes.io/function: |
      exec:
        path: helmize
spec:
  chart: common-app
  repo:
    name: ardikabs
    url: https://charts.ardikabs.com
  version: 0.4.1
  values:
    - values.yaml
    - values/*.yaml
    - values/**/*.yaml
<<EOF
```

Then, you put above Release yaml to the generators field in the `kustomization.yaml`:

```yaml
cat <<EOF > kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generators:
- release.yaml
EOF
```

To generate the manifest you need to use the following command:
```bash

$ kustomize build --enable-alpha-plugins --enable-exec .

```

## More

For more example, please refer to [examples](./examples) directory.
