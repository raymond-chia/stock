FROM python:3.13

WORKDIR /app

COPY pyproject.toml poetry.lock ./
# https://ta-lib.org/install/#linux-debian-packages
RUN wget https://github.com/ta-lib/ta-lib/releases/download/v0.6.4/ta-lib_0.6.4_amd64.deb \
    && dpkg -i ta-lib_0.6.4_amd64.deb \
    && pip install poetry \
    && poetry config virtualenvs.create false \
    && poetry install

CMD ["bash"]
