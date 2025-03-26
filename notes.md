after completion

1. create a dockerfile
2. create a helm chart
3. deploy that on k8s cluster

- also implment go validator



I am using go-blueprint for creating the project structure. you can create go structure from here - https://github.com/golang-standards/project-layout































how I am thinking is:

I take the link from user, make GET req to that link then bring all the responce to /tmp/helm-chart-abc123/chart.tgz 