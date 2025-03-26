after completion

1. create a dockerfile
2. create a helm chart
3. deploy that on k8s cluster

- also implment go validator



I am using go-blueprint for creating the project structure. you can create go structure from here - https://github.com/golang-standards/project-layout


we will try using helm go sdk in cli but here we simply did it using helm pull command




This curl will work in powershell - 

$body = '{\"chartURL\":\"oci://registry-1.docker.io/bitnamicharts/redis\"}'
>> curl.exe -X POST http://localhost:8080/scan -H "Content-Type: application/json" -d "$body"



docker exec helm-scan curl -X POST http://localhost:8080/your-endpoint  -d '{"chartURL": "oci://registry-1.docker.io/bitnamicharts/redis"}'
























how I am thinking is:

I take the link from user, make GET req to that link then bring all the responce to /tmp/helm-chart-abc123/chart.tgz 