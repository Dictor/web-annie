# web-annie
**web-annie** is web interface for annie built with Golang. It can be used to simple video download manager with web interface.

## Installation
### Docker (recommended)
With docker:
```
sudo docker run -p 8080:80 \
                -v /path/to/download/directory:/web-annie/download \
                -v /path/to/config.yaml:/web-annie/config.yaml \
                kimdictor/web-annie
```  
With docker-compose:
```
version: "3"
  services: 
    webannie:
      image: imdictor/web-annie
      ports:
        - "8080:80"
      volumes:
        - /path/to/download/directory:/web-annie/download
        - /path/to/config.yaml:/web-annie/config.yaml
```

### Binary
Download latest [web-annie binary file](https://github.com/Dictor/web-annie/releases) and [annie binary file](https://github.com/iawia002/annie/releases) 
and Place both in same directory ([ffmpeg](https://github.com/iawia002/annie#prerequisites) is required too). 
You can execute web-annie via: 
```
wget "https://github.com/iawia002/annie/releases/download/0.10.3/annie_0.10.3_Linux_64-bit.tar.gz" # download properly binary on your env!
tar -xvfz annie_*
./web-annie
```
