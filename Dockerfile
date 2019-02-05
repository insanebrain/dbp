FROM docker:git
MAINTAINER  Alexandre HAAG <alexandre.haag90@gmail.com>

COPY entrypoint.sh /

RUN chmod a+x entrypoint.sh

COPY dbp /usr/local/bin


RUN mkdir /build

WORKDIR /build

ENTRYPOINT ["/entrypoint.sh"]
CMD ["dbp"]
