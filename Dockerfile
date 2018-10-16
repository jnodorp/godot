FROM busybox
COPY ./godot /godot
CMD ["/godot"]
