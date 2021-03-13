# CI_COMMIT_TAG  [ [source] ](https://github.com/helmwave/helmwave/tree/main/docs/examples/CI_COMMIT_TAG)

Project Structure

```
.
├── README.md
├── helmwave.yml.tpl
└── values.yml

```

```yaml
{% include_relative helmwave.yml.tpl %}
```

```yaml
{%- capture debugging-doc -%}{% include_relative values.yml %}{%- endcapture -%}
{{ debugging-doc | markdownify }}
```

