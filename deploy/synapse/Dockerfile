FROM matrixdotorg/synapse:v1.164.1
COPY homeserver.yaml appservice.yaml log.config render.py /prisme/
ENTRYPOINT ["python3", "/prisme/render.py"]
