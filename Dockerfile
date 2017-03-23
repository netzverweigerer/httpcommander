FROM scratch
ADD httpcommander /
ADD config-examples/test.conf /
ENTRYPOINT ["/httpcommander"]
CMD ["/test.conf"]
