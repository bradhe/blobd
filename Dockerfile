FROM bitnami/minideb:stretch
ADD ./blobd /
RUN install_packages ca-certificates 
CMD ["/blobd"]
