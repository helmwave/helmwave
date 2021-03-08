project: my-project
version: 0.8.4

repositories:
  - name: stable
    url: https://kubernetes-charts.storage.googleapis.com
  - name: bitnami
    url: https://charts.bitnami.com/bitnami

releases:
  - name: nginx
    chart: bitnami/nginx
    options:
      install: true
      namespace: test-nginx
