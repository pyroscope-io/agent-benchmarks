FROM mcr.microsoft.com/dotnet/sdk:6.0


# RUN apt-get update && apt-get -y install time

WORKDIR /opt/agent


ADD Datadog.Linux.ApiWrapper.x64.so .
ADD Datadog.Profiler.Native.so .
COPY fib.csproj fib.cs run.sh ./

RUN dotnet publish --configuration Release -o .

CMD ["/bin/bash", "run.sh"]
