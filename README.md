# web-annie
**web-annie** is web interface for annie built with Golang. It can be used to simple video download manager with web interface.    

![demo](demo_video.gif?raw=true)

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
      image: kimdictor/web-annie
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
tar -xvf annie_*
wget "https://github.com/Dictor/web-annie/releases/download/v1.1.2/web_annie-v1.1.2-linux-amd64.tar.gz" download latest and properly binary on your env!
tar -xvf web_annie*
./web-annie
```

## Configuration
Config file have to be named `config.yaml` and be placed same directory with web-annie binary.  
Full example of `config.yaml`:
```
http_proxy: true
http_proxy_address: 127.0.0.1:9000
download_path: ./download
listen_address: ":80"
ignore_exit_error: false
```
If `config.yaml` doesn't exist, web-annie automatically use default setting. Default setting same as below:
```
http_proxy: false
download_path: ./download
listen_address: ":80"
ignore_exit_error: false
```
