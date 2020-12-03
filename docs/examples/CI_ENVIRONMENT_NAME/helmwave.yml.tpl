project: my-project
version: 0.1.6

releases:
  - name: my-release
    chart: my-chart-repo/{{ env "CI_PROJECT_NAME" }}
    values:
      # Default
      - values/_.yml
      # For specific ENVIRONMENT
      - values/{{ env "CI_ENVIRONMENT_NAME" }}.yml
    options:
      install: true
      namespace: {{ env "CI_ENVIRONMENT_NAME" }}
