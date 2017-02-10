FROM scratch

ADD dist/botella_*_linux-amd64 /botella

ENTRYPOINT ["./botella"]
