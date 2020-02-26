# Kpack Builder Backend

### Prerequisites
* `kpack` installed and configured
* ServiceAccount for the pod that has API access to secrets and kpack images. 

An example Role for the ServiceAccount:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: kpack-builder
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["build.pivotal.io"] #
  resources: ["images"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kpack-builder-binding
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kpack-builder
subjects:
- kind: ServiceAccount
  name: default
  namespace: default
```


To use, run in cluster, expose port 8080. If you want security, supply basic credentials via a secret in the same namespace called `kpack-build-backend`
that looks like this:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: kpack-build-backend
type: Opaque
data:
  username: YWRtaW4=
  password: c2VjcmV0IQ==
```


you make calls to the api with basic auth using:

`/image`
`GET` will return an image list

`/image/username/repo`
`GET` will retrive the details of the image of the name you supply via json like this:
`POST` will allow you to create a new image or updating an existing one by posting json in this format:
`DELETE` will delete the image config

```json
{
  "repo": "https://github.com/someuser/somerepo.git",
  "revision": "master"
}
```
