FROM busybox

LABEL maintainer="julian.schlichtholz@gmail.com"

COPY ./godot /godot
CMD ["/godot"]
