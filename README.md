# Mr. Burns - Micro services testing framework
### Runing Mr. Burns inside Docker container
If you choose to use Mr. Burns Docker [image](https://hub.docker.com/r/gaiaadm/mr-burns/) on Linux, use the following command:
```
docker run -d --name mr-burns -v "$PWD":/go/src/github.com/gaia-adm/mr-burns -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp --link mr-burns-distributor:distributor-link gaiaadm/mr-burns
```
