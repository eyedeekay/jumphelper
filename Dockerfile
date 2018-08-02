FROM alpine:3.7
RUN apk update && apk add go go-tools make musl-dev musl musl-utils git
RUN adduser -g i2pd -D i2pd
#RUN git clone https://github.com/eyedeekay/jumphelper /opt/work
COPY . /opt/work
WORKDIR /opt/work
RUN go get -u github.com/eyedeekay/jumphelper/src
RUN go get -u golang.org/x/time/rate
RUN make deps server
COPY misc/addresses.csv /var/lib/i2pd/addressbook/addresses.csv
RUN chown -R i2pd:i2pd /var/lib/i2pd/addressbook/addresses.csv /opt/work
USER i2pd
VOLUME /opt/work
CMD ./bin/jumphelper -host="0.0.0.0" \
    -share=true \
    -i2p=true \
    -tunname="sam-jumpkelper" \
    -port="7854" \
    -samhost=sam-host \
    -samport="7656" \
    -hostfile=/var/lib/i2pd/addressbook/addresses.csv \
    -subs "http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt"
