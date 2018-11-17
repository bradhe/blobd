FROM bitnami/minideb:stretch
RUN install_packages ca-certificates 
ADD ./blobd /
CMD ["/blobd"]
