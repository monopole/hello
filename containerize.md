# Containerize Hello

### Bash functions for container building

Define these first; they get used below.

Requirments: _Go_ compiler, _curl_, and _docker_.

<!-- @funcBuildVersionedExecutable -->
```sh
function buildVersionedExecutable {
  local tmpDir=$1
  local githubUser=$2
  local pgmName=$3
  local version=$4

  local package=github.com/${githubUser}/${pgmName}
  local newPgm=$tmpDir/${pgmName}_${version}

  GOPATH=$tmpDir go get -d $package

  cat $tmpDir/src/${package}/${pgmName}.go |\
      sed 's/version = 0/version = '${version}'/' \
      >${newPgm}.go

  echo Compiling ${newPgm}.go
  GOPATH=$tmpDir CGO_ENABLED=0 GOOS=linux go build \
      -o $tmpDir/${pgmName} \
      -a -installsuffix cgo ${newPgm}.go
}
```

<!-- @funcRunAndQuitRawBinaryToTest -->
```sh
function runAndQuitRawBinaryToTest {
  local tmpDir=$1
  local pgmName=$2
  local port=$3

  echo Running server $tmpDir/$pgmName
  ALT_GREETING=salutations \
      $tmpDir/$pgmName --enableRiskyFeature --port $port &

  # Let it get ready
  sleep 2

  # Dump html to stdout
  curl --fail --silent -m 1 localhost:$port/godzilla

  # Send query of death
  curl --fail --silent -m 1 localhost:$port/quit
  echo Server stopped
}
```

<!-- @funcBuildDockerImage -->
```sh
function buildDockerImage {
  local tmpDir=$1
  local pgmName=$2
  local version=$3

  # Repo holds just one image, give repo same name as image.
  local dockerRepo=$pgmName

  local dockerFile=$tmpDir/Dockerfile
  cat <<EOF >$dockerFile
FROM scratch
ADD $pgmName /
CMD ["/$pgmName"]
EOF
  echo Docker build
  docker build -t $dockerRepo:$version -f $dockerFile $tmpDir
  echo End docker build
}

```


<!-- @funcRunAndQuitInsideDockerToTest -->
```sh
function runAndQuitInsideDockerToTest {
  local pgmName=$1
  local version=$2
  local port=$3

  echo Docker run, mapping $port to internal 8080
  docker run -d -p $port:8080 $pgmName:$version
  sleep 3
  docker ps | grep $pgmName

  echo Requesting docker server
  curl -m 1 localhost:$port/kingGhidorah
  curl -m 1 localhost:$port/quit
}
```

<!-- @funcPushToDockerHub -->
```sh
function pushToDockerHub {
  local dockerUser=$1
  local pgmName=$2
  local version=$3

  local repoName=$pgmName

  local id=$(docker images |\
      grep $pgmName | grep " $version " | awk '{printf $3}')
  docker tag $id $dockerUser/$repoName:$version
  docker push $dockerUser/$repoName:$version
}
```

<!-- @funcBuildContainer -->
```sh
function buildContainer {
  local githubOrg=$1
  local pgmName=$2
  local version=$3
  local testPort=$4
  local tmpDir=$(mktemp -d)

  echo tmpDir=$tmpDir
  buildVersionedExecutable $tmpDir $githubOrg $pgmName $version
  runAndQuitRawBinaryToTest $tmpDir $pgmName $testPort

  buildDockerImage $tmpDir $pgmName $version
  docker images --no-trunc | grep $pgmName
  sleep 4
  runAndQuitInsideDockerToTest $pgmName $version $testPort
}
```

<!-- @funcRemoveLocalImage -->
```sh
function removeLocalImage {
  local pgmName=$1
  local version=$2

  echo docker rmi $pgmName:$version
  docker rmi $pgmName:$version
  id=$(docker images | grep $pgmName | grep " $version " | awk '{printf $3}')
  echo docker rmi -f $id
  docker rmi -f $id
}
```

### Create and upload images to [hub.docker.com](https://hub.docker.com/r/monopole/hello)

<!-- @setUp -->
```sh
dockerUser=monopole
githubOrg=monopole
```

<!-- @login -->
```sh
printf "\nEnter docker password, followed by C-d: "
docker login --username=$dockerUser --password-stdin
```

<!-- @doVersion1 -->
```sh
buildContainer $githubOrg hello 1 8999
pushToDockerHub $dockerUser hello 1
```

<!-- @doVersion2 -->
```sh
buildContainer $githubOrg hello 2 8999
pushToDockerHub $dockerUser hello 2
```

### Test images

Remove the images from the local cache, then run them.
This forces a new pull.

```sh
removeLocalImage hello 1
removeLocalImage hello 2
```

```sh
docker run -d -p 8999:8080 docker.io/$dockerUser/hello:1
curl -m 1 localhost:8999/shouldBeV1
curl -m 1 localhost:8999/quit
```

```sh
docker run -d -p 8999:8080 docker.io/$dockerUser/hello:2
curl -m 1 localhost:8999/shouldBeV2
curl -m 1 localhost:8999/quit
```
