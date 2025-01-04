docker build -t treasure-doc .
docker save -o treasure-doc.tar.gz treasure-doc
sudo chmod 777 treasure-doc.tar.gz