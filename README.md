# Helmize

> A KRM function to enable Helm on Kustomize with Glob support

The motivation for this project is quite simple. While [Helm integration](https://kubectl.docs.kubernetes.io/references/kustomize/builtins/#_helmchartinflationgenerator_) in Kustomize is already available, it falls short in scenarios where dynamic use of Glob is required to fetch all necessary values.yaml files. In such cases, it lacks support.
With the likelihood of [glob support](https://github.com/kubernetes-sigs/kustomize/issues/119) being added to Kustomize remaining uncertain for the foreseeable future, this project was initiated to address this limitation.

## Usage

The Release specification:

```bash
cat <<EOF > release-simple.yaml
apiVersion: toolkit.ardikabs.com/v1alpha1
kind: HelmRelease
metadata:
  name: simple-a
  namespace: default
  annotations:
    config.kubernetes.io/function: |
      exec:
        # path: helmize
        path: ../../bin/helmize
spec:
  chart: minecraft
  repo:
    name: minecraft
    url: https://itzg.github.io/minecraft-server-charts
EOF

cat <<EOF > release-with-glob.yaml
apiVersion: toolkit.ardikabs.com/v1alpha1
kind: HelmRelease
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
EOF

cat <<EOF > release-with-oci-repo.yaml
apiVersion: toolkit.ardikabs.com/v1alpha1
kind: HelmRelease
metadata:
  name: envoy-gateway
  namespace: envoy-gateway-system
  annotations:
    config.kubernetes.io/function: |
      exec:
        # path: helmize
        path: ../../bin/helmize
spec:
  repo:
    url: oci://docker.io/envoyproxy/gateway-helm
  version: v0.5.0
  includeCRDs: true
  createNamespace: true
EOF
```

Then, you put above Release yaml to the generators field in the `kustomization.yaml`:

```bash
cat <<EOF > kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generators:
- release-simple.yaml
- release-with-glob.yaml
- release-with-oci.yaml
EOF
```

To generate the manifest you need to use the following command:
```bash

$ kustomize build --enable-alpha-plugins --enable-exec .

```

## More

For more example, please refer to [examples](./examples) directory.
