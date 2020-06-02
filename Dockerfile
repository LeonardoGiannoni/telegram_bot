FROM archlinux
WORKDIR /app
COPY run.sh .
COPY labprogettazione_telegrambot .
EXPOSE 8081
ENTRYPOINT ["/app/run.sh"]
