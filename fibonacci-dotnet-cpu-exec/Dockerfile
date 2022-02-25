FROM mcr.microsoft.com/dotnet/sdk:6.0

WORKDIR /opt/agent

COPY --from=pyroscope/pyroscope:latest /usr/bin/pyroscope /usr/bin/pyroscope
COPY fib.csproj fib.cs run.sh ./

RUN dotnet publish --configuration Release -o .

CMD ["/bin/bash", "run.sh"]
