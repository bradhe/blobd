FROM bitnami/minideb:stretch
ADD ./blobd /
CMD ["/blobd"]
