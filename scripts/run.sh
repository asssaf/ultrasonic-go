: ${IMAGE_NAME:=asssaf/ultrasonic:latest}
docker run --rm -it --privileged --device /dev/ttyAMA0 "$IMAGE_NAME" $*
