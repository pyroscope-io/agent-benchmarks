FROM mcr.microsoft.com/dotnet/sdk:6.0


WORKDIR /opt/agent


RUN wget https://github.com/pyroscope-io/pyroscope-dotnet/releases/download/v0.02-pyroscope/Datadog.Linux.ApiWrapper.x64.so
RUN wget https://github.com/pyroscope-io/pyroscope-dotnet/releases/download/v0.02-pyroscope/Datadog.Profiler.Native.so


COPY fib.csproj fib.cs run.sh ./

RUN dotnet publish --configuration Release -o .

CMD ["/bin/bash", "run.sh"]
