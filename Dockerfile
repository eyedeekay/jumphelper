FROM alpine:3.7
RUN apk update && apk add go go-tools make musl-dev musl musl-utils git
RUN adduser -g i2pd -D i2pd
#RUN git clone https://github.com/eyedeekay/jumphelper /opt/work
COPY . /opt/work
WORKDIR /opt/work
RUN go get -u github.com/eyedeekay/jumphelper/src
RUN go get -u golang.org/x/time/rate
RUN make server
COPY addresses.csv /var/lib/i2pd/addressbook/addresses.csv
RUN chown i2pd:i2pd /var/lib/i2pd/addressbook/addresses.csv
USER i2pd
CMD ./bin/jumphelper -hostfile=/var/lib/i2pd/addressbook/addresses.csv
