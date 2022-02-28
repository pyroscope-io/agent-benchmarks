FROM openjdk:16-slim-bullseye
WORKDIR /opt/agent

RUN apt-get update && apt-get install ca-certificates -y && \
    update-ca-certificates && apt-get install -y curl git unzip zip
RUN git clone https://github.com/pyroscope-io/pyroscope-java.git
RUN cd pyroscope-java && \
    ./gradlew shadowJar && \
    cp agent/build/libs/pyroscope.jar ../pyroscope.jar

RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN curl -s https://get.sdkman.io | bash && source ~/.sdkman/bin/sdkman-init.sh && sdk install kotlin
COPY fib.kt run.sh ./
RUN source ~/.sdkman/bin/sdkman-init.sh && kotlinc fib.kt -include-runtime -d fib.jar

CMD ["/bin/bash", "run.sh"]
