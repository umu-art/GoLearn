FROM elasticsearch:8.14.3
LABEL authors="vikazeni"

USER 0:0
COPY ./homework_4/iac/image/go-homework-4-elastic/config /usr/share/elasticsearch/config
RUN chmod 777 -R /usr/share/elasticsearch/config
USER 1000:0
