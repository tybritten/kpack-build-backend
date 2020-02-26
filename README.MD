# Kpack Builder Backend

*Requires kpack installed and configured*

To use, run in cluster, expose port 8080. supply basic credentials via a secret in the same namespace called `kpack-build-backend`
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
