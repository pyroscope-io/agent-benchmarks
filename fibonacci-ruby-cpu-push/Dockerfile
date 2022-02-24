FROM ruby:3.0.2

WORKDIR /opt/agent
COPY Gemfile Gemfile.lock fib.rb ./
RUN bundle install

CMD [ "ruby", "fib.rb" ]
