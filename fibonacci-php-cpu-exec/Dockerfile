FROM php:8.1.3

WORKDIR /opt/agent

COPY --from=pyroscope/pyroscope:latest /usr/bin/pyroscope /usr/bin/pyroscope
COPY fib.php run.sh ./

CMD ["/bin/bash", "run.sh"]
