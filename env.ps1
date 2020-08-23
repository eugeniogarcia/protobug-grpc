$ubicacion=pwd
$env:GOPATH=$env:GOPATH+";"+$ubicacion
# Con esto cuando hagamos go install xxxxxxxxx dejara el resultado en el directorio bin dentro del pwd
$env:GOBIB = $ubicacion

