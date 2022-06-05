# cos-tools

requires:
- [jq](https://stedolan.github.io/jq/)
- curl
- [stern](https://github.com/stern/stern)
- [ocm](https://github.com/openshift-online/ocm-cli)
- kubectl

define the following env vars:
- COS_BASE_PATH -> base URL for the managed conenctor service
- KAS_BASE_PATH -> base URL for the managed kafka service


