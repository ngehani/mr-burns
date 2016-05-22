# Mr. Burns - Micro services testing framework
### Running Mr. Burns inside Docker container
If you choose to use Mr. Burns Docker [image](https://hub.docker.com/r/gaiaadm/mr-burns/) on Linux, use the following command:
```bash
docker run -d --name burns --log-driver=json-file -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp gaiaadm/mr-burns
```
