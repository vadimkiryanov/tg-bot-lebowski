# tg-bot-lebowski
docker build -t tg-bot-lebowski .                
docker run --name tg-bot-lebowski -p 80:80 --env-file .env tg-bot-lebowski
docker run -d --name tg-bot-lebowski -p 80:80 --env-file .env --restart=unless-stopped tg-bot-lebowski
docker rmi -f tg-bot-lebowski

sudo apt-get install build-essential 
