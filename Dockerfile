FROM python:3.13-alpine

WORKDIR /usr/src/app
RUN mkdir output

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY *.py ./

ENTRYPOINT [ "python", "contrast_application_licenses.py" ]
