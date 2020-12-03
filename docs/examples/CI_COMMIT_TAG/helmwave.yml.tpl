project: my-project # Имя проекта
version: 0.1.6 # Версия helmwave

releases:
  - name: my-release
    chart: my-chart-repo/my-app
    values:
      - values.yml
    options:
      install: true
      namespace: my-namespace
