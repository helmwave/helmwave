project: {{ env "PROJECT_NAME" | default "bad" }}

repositories:
- name: stable
  url: https://kubernetes-charts.storage.googleapis.com
- name: bitnami
  url: https://charts.bitnami.com/bitnami


releases:
- name: nginx
  chart:
    name: bitnami/nginx
  namespace: {{ requiredEnv "NAMESPACE" }}
