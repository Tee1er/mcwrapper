cd data/server/
datestamp=$(date +"%Y-%m-%d_%H-%M-%S")
zip -r "../FW8 backup $datestamp" ./worlds/FW8
