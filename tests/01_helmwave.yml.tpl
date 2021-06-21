project: {{ env "PROJECT_NAME" | default "bad" }}

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
    namespace: {{ requiredEnv "NAMESPACE" }}
