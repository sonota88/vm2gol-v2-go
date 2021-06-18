素朴な自作言語のコンパイラをGoに移植した - memo88  
https://memo88.hatenablog.com/entry/2020/09/25/073200

```sh
## Build Docker image

docker build \
  --build-arg USER=$USER \
  --build-arg GROUP=$(id -gn) \
  -t vm2gol-v2:go .

## Run

docker run --rm -it \
  -v"$(pwd):/home/${USER}/work" \
  vm2gol-v2:go /bin/bash
```
