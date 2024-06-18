## first setup

  - 

## system process

 - https://stackoverflow.com/questions/43650183/keep-running-go-server-as-background-process
 - https://github.com/emersion/go-autostart
 - https://github.com/radovskyb/gobeat
 - https://github.com/tillberg/autorestart
 - https://github.com/tillberg/autoinstall

## create new runtime

  - upload Dockerfile
    - get runtime name from request
      - generate a random one if not provided
    - copy it to a folder
    - set KV to runtime name - file name

## create a new function

  - upload zip of function
    - name of runtime
    - tag (name of container)
    - docker settings
      - port
      - volume
  - unzip
  - copy to temp folder with runtime file
  - call build on temp folder
    - `docker build -t $NAME .`
  - run built image
    - set to common network
